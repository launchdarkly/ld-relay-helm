package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	corev1 "k8s.io/api/core/v1"
)

func (s *TemplateTest) TestConfigCanSetCommonLabels() {
	options := &helm.Options{
		SetValues: map[string]string{
			"commonLabels.environment": "production",
			"commonLabels.team":        "platform",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/config.yaml"})
	var config corev1.ConfigMap
	helm.UnmarshalK8SYaml(s.T(), output, &config)

	s.Require().Equal("production", config.Labels["environment"])
	s.Require().Equal("platform", config.Labels["team"])
}
