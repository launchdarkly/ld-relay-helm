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

// The chart adds no writable volume for a read-only root filesystem. Users who
// need one, for example for offline mode, mount it themselves.
func (s *TemplateTest) TestReadOnlyRootFilesystemAddsNoVolumes() {
	deployment := s.renderDeployment(map[string]string{
		"securityContext.readOnlyRootFilesystem": "true",
	})

	s.Require().True(*deployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem)
	s.Require().Len(deployment.Spec.Template.Spec.Volumes, 1)
	s.Require().Empty(deployment.Spec.Template.Spec.Containers[0].VolumeMounts)
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
