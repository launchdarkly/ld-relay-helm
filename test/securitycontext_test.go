package test

import (
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

func (s *TemplateTest) TestDoesNotMountTmpWhenTheUserAlreadyMountsIt() {
	deployment := s.renderDeployment(map[string]string{
		"securityContext.readOnlyRootFilesystem": "true",
		"relay.extraVolumes[0].name":             "my-tmp",
		"relay.extraVolumes[0].emptyDir.medium":  "Memory",
		"relay.extraVolumeMounts[0].name":        "my-tmp",
		"relay.extraVolumeMounts[0].mountPath":   "/tmp/",
	})

	s.Require().Nil(findVolume(deployment.Spec.Template.Spec.Volumes, "ld-relay-tmp"))

	mounts := deployment.Spec.Template.Spec.Containers[0].VolumeMounts
	s.Require().Len(mounts, 1)
	s.Require().Equal("my-tmp", mounts[0].Name)
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
