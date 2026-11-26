package test

import (
	"path/filepath"
	"strings"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (s *TemplateTest) renderDeployment(values map[string]string) appsv1.Deployment {
	options := &helm.Options{
		SetValues:      values,
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	return deployment
}

func (s *TemplateTest) renderDeploymentFromFile(name string) appsv1.Deployment {
	valuesFile, err := filepath.Abs(filepath.Join("testdata", name))
	s.Require().NoError(err)

	options := &helm.Options{
		ValuesFiles:    []string{valuesFile},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	return deployment
}

func (s *TemplateTest) TestContainerDefaultsToARestrictedSecurityContext() {
	deployment := s.renderDeployment(nil)

	sc := deployment.Spec.Template.Spec.Containers[0].SecurityContext
	s.Require().NotNil(sc)

	s.Require().False(*sc.AllowPrivilegeEscalation)
	s.Require().True(*sc.RunAsNonRoot)
	s.Require().Equal([]corev1.Capability{"ALL"}, sc.Capabilities.Drop)
	s.Require().Equal(corev1.SeccompProfileTypeRuntimeDefault, sc.SeccompProfile.Type)
	s.Require().True(*sc.ReadOnlyRootFilesystem)

	// The image sets the user. The chart does not pin it.
	s.Require().Nil(sc.RunAsUser)
	s.Require().Nil(sc.RunAsGroup)
}

func (s *TemplateTest) TestPodDefaultsToARestrictedSecurityContext() {
	deployment := s.renderDeployment(nil)

	sc := deployment.Spec.Template.Spec.SecurityContext
	s.Require().NotNil(sc)

	s.Require().Equal(corev1.SeccompProfileTypeRuntimeDefault, sc.SeccompProfile.Type)

	// runAsNonRoot is set on the relay container only, so root init containers still run.
	s.Require().Nil(sc.RunAsNonRoot)
	s.Require().Nil(sc.RunAsUser)
	s.Require().Nil(sc.RunAsGroup)
	s.Require().Nil(sc.FSGroup)
}

func (s *TemplateTest) TestCanRemoveADefaultSecurityContextField() {
	deployment := s.renderDeployment(map[string]string{
		"securityContext.runAsNonRoot": "null",
	})

	sc := deployment.Spec.Template.Spec.Containers[0].SecurityContext
	s.Require().Nil(sc.RunAsNonRoot)
	s.Require().False(*sc.AllowPrivilegeEscalation)
}

func (s *TemplateTest) TestCanRemoveAllSecurityContextDefaults() {
	deployment := s.renderDeploymentFromFile("securitycontext_null.yaml")

	s.Require().Nil(deployment.Spec.Template.Spec.SecurityContext)
	s.Require().Nil(deployment.Spec.Template.Spec.Containers[0].SecurityContext)
}

// Helm merges maps, so an empty map in a values file keeps every default.
func (s *TemplateTest) TestEmptySecurityContextKeepsTheDefaults() {
	deployment := s.renderDeploymentFromFile("securitycontext_empty.yaml")

	s.Require().Equal(corev1.SeccompProfileTypeRuntimeDefault, deployment.Spec.Template.Spec.SecurityContext.SeccompProfile.Type)

	sc := deployment.Spec.Template.Spec.Containers[0].SecurityContext
	s.Require().True(*sc.RunAsNonRoot)
	s.Require().False(*sc.AllowPrivilegeEscalation)
	s.Require().True(*sc.ReadOnlyRootFilesystem)
}

func (s *TemplateTest) findVolume(deployment appsv1.Deployment, name string) *corev1.Volume {
	for i := range deployment.Spec.Template.Spec.Volumes {
		if deployment.Spec.Template.Spec.Volumes[i].Name == name {
			return &deployment.Spec.Template.Spec.Volumes[i]
		}
	}
	return nil
}

func (s *TemplateTest) mountsAt(deployment appsv1.Deployment, path string) []corev1.VolumeMount {
	var mounts []corev1.VolumeMount
	for _, m := range deployment.Spec.Template.Spec.Containers[0].VolumeMounts {
		if strings.TrimSuffix(m.MountPath, "/") == path {
			mounts = append(mounts, m)
		}
	}
	return mounts
}

// The root filesystem is read-only by default, so the chart mounts a writable
// emptyDir at /tmp. Offline mode expands its archive there.
func (s *TemplateTest) TestMountsAWritableTmpByDefault() {
	deployment := s.renderDeployment(nil)

	s.Require().True(*deployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem)

	volume := s.findVolume(deployment, s.Release+"-tmp")
	s.Require().NotNil(volume)
	s.Require().NotNil(volume.EmptyDir)

	mounts := s.mountsAt(deployment, "/tmp")
	s.Require().Len(mounts, 1)
	s.Require().Equal(s.Release+"-tmp", mounts[0].Name)
	s.Require().False(mounts[0].ReadOnly)
}

func (s *TemplateTest) TestOfflineModeGetsAWritableTmp() {
	deployment := s.renderDeployment(map[string]string{
		"relay.volume.offline": "relay-archive.tar.gz",
		"relay.volume.definition.persistentVolumeClaim.claimName": "ld-relay-pvc",
	})

	s.Require().True(*deployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem)
	s.Require().Len(s.mountsAt(deployment, "/tmp"), 1)
	s.Require().Len(s.mountsAt(deployment, "/mnt/volume"), 1)
}

func (s *TemplateTest) TestCanTurnOffReadOnlyRootFilesystem() {
	deployment := s.renderDeployment(map[string]string{
		"securityContext.readOnlyRootFilesystem": "false",
	})

	sc := deployment.Spec.Template.Spec.Containers[0].SecurityContext
	s.Require().False(*sc.ReadOnlyRootFilesystem)
	s.Require().True(*sc.RunAsNonRoot)
}

// A user who sets one securityContext field keeps the other defaults,
// including the read-only root filesystem.
func (s *TemplateTest) TestOverridingASecurityContextFieldKeepsReadOnlyRootFilesystem() {
	deployment := s.renderDeployment(map[string]string{
		"securityContext.runAsUser": "10001",
	})

	sc := deployment.Spec.Template.Spec.Containers[0].SecurityContext
	s.Require().Equal(int64(10001), *sc.RunAsUser)
	s.Require().True(*sc.ReadOnlyRootFilesystem)
	s.Require().True(*sc.RunAsNonRoot)
	s.Require().False(*sc.AllowPrivilegeEscalation)
	s.Require().Len(s.mountsAt(deployment, "/tmp"), 1)
}

// If the user mounts their own volume at /tmp, the chart does not add a second
// one. Two mounts at the same path would make the pod invalid.
func (s *TemplateTest) TestUserTmpMountReplacesTheChartTmpVolume() {
	for _, path := range []string{"/tmp", "/tmp/"} {
		deployment := s.renderDeployment(map[string]string{
			"relay.extraVolumes[0].name":               "tmp",
			"relay.extraVolumes[0].emptyDir.medium":    "Memory",
			"relay.extraVolumes[0].emptyDir.sizeLimit": "64Mi",
			"relay.extraVolumeMounts[0].name":          "tmp",
			"relay.extraVolumeMounts[0].mountPath":     path,
		})

		s.Require().Nil(s.findVolume(deployment, s.Release+"-tmp"), path)
		s.Require().NotNil(s.findVolume(deployment, "tmp"), path)

		mounts := s.mountsAt(deployment, "/tmp")
		s.Require().Len(mounts, 1, path)
		s.Require().Equal("tmp", mounts[0].Name, path)
	}
}

func (s *TemplateTest) TestDoesNotMountTheServiceAccountTokenByDefault() {
	deployment := s.renderDeployment(nil)

	s.Require().NotNil(deployment.Spec.Template.Spec.AutomountServiceAccountToken)
	s.Require().False(*deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}

func (s *TemplateTest) TestCanLeaveAutomountServiceAccountTokenUnset() {
	deployment := s.renderDeployment(map[string]string{
		"automountServiceAccountToken": "null",
	})

	s.Require().Nil(deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}

func (s *TemplateTest) TestCanEnableAutomountingTheServiceAccountToken() {
	deployment := s.renderDeployment(map[string]string{
		"automountServiceAccountToken": "true",
	})

	s.Require().NotNil(deployment.Spec.Template.Spec.AutomountServiceAccountToken)
	s.Require().True(*deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}
