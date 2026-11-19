package test

import (
	"path/filepath"

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

func findVolume(volumes []corev1.Volume, name string) *corev1.Volume {
	for i := range volumes {
		if volumes[i].Name == name {
			return &volumes[i]
		}
	}
	return nil
}

func (s *TemplateTest) TestContainerDefaultsToARestrictedSecurityContext() {
	deployment := s.renderDeployment(nil)

	sc := deployment.Spec.Template.Spec.Containers[0].SecurityContext
	s.Require().NotNil(sc)

	s.Require().False(*sc.AllowPrivilegeEscalation)
	s.Require().True(*sc.RunAsNonRoot)
	s.Require().Equal([]corev1.Capability{"ALL"}, sc.Capabilities.Drop)
	s.Require().Equal(corev1.SeccompProfileTypeRuntimeDefault, sc.SeccompProfile.Type)

	// The image sets the user. The chart does not pin it.
	s.Require().Nil(sc.RunAsUser)
	s.Require().Nil(sc.RunAsGroup)
	s.Require().Nil(sc.ReadOnlyRootFilesystem)
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
}

func (s *TemplateTest) TestDoesNotMountTmpByDefault() {
	deployment := s.renderDeployment(nil)

	s.Require().Nil(findVolume(deployment.Spec.Template.Spec.Volumes, "ld-relay-tmp"))
	s.Require().Empty(deployment.Spec.Template.Spec.Containers[0].VolumeMounts)
}

// A read-only root filesystem needs writable scratch space. Offline mode
// expands its archive under TMPDIR.
func (s *TemplateTest) TestMountsTmpWithAReadOnlyRootFilesystem() {
	deployment := s.renderDeployment(map[string]string{
		"securityContext.readOnlyRootFilesystem": "true",
	})

	tmpVolume := findVolume(deployment.Spec.Template.Spec.Volumes, "ld-relay-tmp")
	s.Require().NotNil(tmpVolume, "expected an ld-relay-tmp volume")
	s.Require().NotNil(tmpVolume.EmptyDir)

	mounts := deployment.Spec.Template.Spec.Containers[0].VolumeMounts
	s.Require().Len(mounts, 1)
	s.Require().Equal("ld-relay-tmp", mounts[0].Name)
	s.Require().Equal("/tmp", mounts[0].MountPath)
}

func (s *TemplateTest) TestDoesNotMountTmpWhenReadOnlyRootFilesystemIsFalse() {
	deployment := s.renderDeployment(map[string]string{
		"securityContext.readOnlyRootFilesystem": "false",
	})

	s.Require().Nil(findVolume(deployment.Spec.Template.Spec.Volumes, "ld-relay-tmp"))
	s.Require().Empty(deployment.Spec.Template.Spec.Containers[0].VolumeMounts)
}

func (s *TemplateTest) TestMountsTmpAfterTheOtherVolumes() {
	deployment := s.renderDeployment(map[string]string{
		"securityContext.readOnlyRootFilesystem":  "true",
		"relay.volume.definition.emptyDir.medium": "Memory",
		"relay.secrets[0].volumePath":             "my-secret-path",
		"relay.secrets[0].secretKey":              "ld-relay",
		"relay.secrets[0].secretName":             "sdk-key",
		"relay.extraVolumes[0].name":              "data",
		"relay.extraVolumes[0].emptyDir.medium":   "Memory",
		"relay.extraVolumeMounts[0].name":         "data",
		"relay.extraVolumeMounts[0].mountPath":    "/data",
	})

	volumes := deployment.Spec.Template.Spec.Volumes
	s.Require().Len(volumes, 5)
	s.Require().Equal("ld-relay-tmp", volumes[len(volumes)-1].Name)

	mounts := deployment.Spec.Template.Spec.Containers[0].VolumeMounts
	s.Require().Len(mounts, 4)
	s.Require().Equal("/mnt/volume", mounts[0].MountPath)
	s.Require().Equal("/mnt/secrets/", mounts[1].MountPath)
	s.Require().Equal("/data", mounts[2].MountPath)
	s.Require().Equal("ld-relay-tmp", mounts[3].Name)
	s.Require().Equal("/tmp", mounts[3].MountPath)
}

func (s *TemplateTest) TestDoesNotMountTmpWhenTheUserAlreadyMountsIt() {
	for _, mountPath := range []string{"/tmp", "/tmp/", "/tmp//", "/./tmp"} {
		deployment := s.renderDeployment(map[string]string{
			"securityContext.readOnlyRootFilesystem": "true",
			"relay.extraVolumes[0].name":             "my-tmp",
			"relay.extraVolumes[0].emptyDir.medium":  "Memory",
			"relay.extraVolumeMounts[0].name":        "my-tmp",
			"relay.extraVolumeMounts[0].mountPath":   mountPath,
		})

		s.Require().Nil(findVolume(deployment.Spec.Template.Spec.Volumes, "ld-relay-tmp"), mountPath)

		mounts := deployment.Spec.Template.Spec.Containers[0].VolumeMounts
		s.Require().Len(mounts, 1, mountPath)
		s.Require().Equal("my-tmp", mounts[0].Name, mountPath)
	}
}

func (s *TemplateTest) TestLeavesAutomountServiceAccountTokenUnsetByDefault() {
	deployment := s.renderDeployment(nil)

	s.Require().Nil(deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}

func (s *TemplateTest) TestCanDisableAutomountingTheServiceAccountToken() {
	deployment := s.renderDeployment(map[string]string{
		"automountServiceAccountToken": "false",
	})

	s.Require().NotNil(deployment.Spec.Template.Spec.AutomountServiceAccountToken)
	s.Require().False(*deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}

func (s *TemplateTest) TestCanEnableAutomountingTheServiceAccountToken() {
	deployment := s.renderDeployment(map[string]string{
		"automountServiceAccountToken": "true",
	})

	s.Require().NotNil(deployment.Spec.Template.Spec.AutomountServiceAccountToken)
	s.Require().True(*deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}
