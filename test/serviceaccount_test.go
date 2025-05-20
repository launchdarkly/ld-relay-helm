package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	corev1 "k8s.io/api/core/v1"
)

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
