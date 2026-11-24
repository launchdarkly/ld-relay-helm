package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	corev1 "k8s.io/api/core/v1"
)

func (s *TemplateTest) TestServiceAccountNamespace() {
	s.assertTemplateNamespaces("ServiceAccount", "templates/serviceaccount.yaml", nil)
}

func (s *TemplateTest) TestServiceAccountCanSetCommonLabels() {
	options := &helm.Options{
		SetValues: map[string]string{
			"serviceAccount.create":    "true",
			"commonLabels.environment": "production",
			"commonLabels.team":        "platform",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/serviceaccount.yaml"})
	var serviceAccount corev1.ServiceAccount
	helm.UnmarshalK8SYaml(s.T(), output, &serviceAccount)

	s.Require().Equal("production", serviceAccount.Labels["environment"])
	s.Require().Equal("platform", serviceAccount.Labels["team"])
}

func (s *TemplateTest) renderServiceAccount(values map[string]string) corev1.ServiceAccount {
	options := &helm.Options{
		SetValues:      values,
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/serviceaccount.yaml"})
	var serviceAccount corev1.ServiceAccount
	helm.UnmarshalK8SYaml(s.T(), output, &serviceAccount)

	return serviceAccount
}

func (s *TemplateTest) TestServiceAccountDoesNotMountTheTokenByDefault() {
	serviceAccount := s.renderServiceAccount(nil)

	s.Require().NotNil(serviceAccount.AutomountServiceAccountToken)
	s.Require().False(*serviceAccount.AutomountServiceAccountToken)
}

func (s *TemplateTest) TestServiceAccountCanMountTheToken() {
	serviceAccount := s.renderServiceAccount(map[string]string{
		"automountServiceAccountToken": "true",
	})

	s.Require().NotNil(serviceAccount.AutomountServiceAccountToken)
	s.Require().True(*serviceAccount.AutomountServiceAccountToken)
}

func (s *TemplateTest) TestServiceAccountCanLeaveAutomountServiceAccountTokenUnset() {
	serviceAccount := s.renderServiceAccount(map[string]string{
		"automountServiceAccountToken": "null",
	})

	s.Require().Nil(serviceAccount.AutomountServiceAccountToken)
}
