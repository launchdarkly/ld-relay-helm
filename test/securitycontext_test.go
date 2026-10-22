package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (s *TemplateTest) renderDefaultDeployment() appsv1.Deployment {
	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	return deployment
}

func (s *TemplateTest) TestContainerDefaultsToARestrictedSecurityContext() {
	deployment := s.renderDefaultDeployment()

	sc := deployment.Spec.Template.Spec.Containers[0].SecurityContext
	s.Require().NotNil(sc)

	s.Require().False(*sc.AllowPrivilegeEscalation)
	s.Require().True(*sc.ReadOnlyRootFilesystem)
	s.Require().True(*sc.RunAsNonRoot)
	s.Require().EqualValues(1000, *sc.RunAsUser)
	s.Require().EqualValues(1000, *sc.RunAsGroup)
	s.Require().Equal([]corev1.Capability{"ALL"}, sc.Capabilities.Drop)
	s.Require().Equal(corev1.SeccompProfileTypeRuntimeDefault, sc.SeccompProfile.Type)
}

func (s *TemplateTest) TestPodDefaultsToARestrictedSecurityContext() {
	deployment := s.renderDefaultDeployment()

	sc := deployment.Spec.Template.Spec.SecurityContext
	s.Require().NotNil(sc)

	s.Require().True(*sc.RunAsNonRoot)
	s.Require().EqualValues(1000, *sc.RunAsUser)
	s.Require().EqualValues(1000, *sc.RunAsGroup)
	s.Require().EqualValues(1000, *sc.FSGroup)
	s.Require().Equal(corev1.SeccompProfileTypeRuntimeDefault, sc.SeccompProfile.Type)
}

// The container defaults to a read-only root filesystem, so it needs writable
// scratch space: offline mode expands its archive under TMPDIR.
func (s *TemplateTest) TestAlwaysMountsAWritableTmpVolume() {
	deployment := s.renderDefaultDeployment()

	var tmpVolume *corev1.Volume
	for i, v := range deployment.Spec.Template.Spec.Volumes {
		if v.Name == "ld-relay-tmp" {
			tmpVolume = &deployment.Spec.Template.Spec.Volumes[i]
		}
	}
	s.Require().NotNil(tmpVolume, "expected an ld-relay-tmp volume")
	s.Require().NotNil(tmpVolume.EmptyDir)

	mounts := deployment.Spec.Template.Spec.Containers[0].VolumeMounts
	s.Require().Equal("ld-relay-tmp", mounts[len(mounts)-1].Name)
	s.Require().Equal("/tmp", mounts[len(mounts)-1].MountPath)
}

func (s *TemplateTest) TestDoesNotAutomountTheServiceAccountTokenByDefault() {
	deployment := s.renderDefaultDeployment()

	s.Require().NotNil(deployment.Spec.Template.Spec.AutomountServiceAccountToken)
	s.Require().False(*deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}

func (s *TemplateTest) TestCanOptInToAutomountingTheServiceAccountToken() {
	options := &helm.Options{
		SetValues: map[string]string{
			"automountServiceAccountToken": "true",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().NotNil(deployment.Spec.Template.Spec.AutomountServiceAccountToken)
	s.Require().True(*deployment.Spec.Template.Spec.AutomountServiceAccountToken)
}
