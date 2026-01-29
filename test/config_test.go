package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	corev1 "k8s.io/api/core/v1"
)

func (s *TemplateTest) TestConfigNameRespectsFullnameOverride() {
	options := &helm.Options{
		SetValues: map[string]string{
			"fullnameOverride": "my-custom-relay",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/config.yaml"})
	var config corev1.ConfigMap
	helm.UnmarshalK8SYaml(s.T(), output, &config)

	s.Require().Equal("my-custom-relay-config", config.Name)
}

func (s *TemplateTest) TestConfigNameIncludesReleaseNameByDefault() {
	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/config.yaml"})
	var config corev1.ConfigMap
	helm.UnmarshalK8SYaml(s.T(), output, &config)

	s.Require().Contains(config.Name, s.Release)
}

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
