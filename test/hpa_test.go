package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
)

func (s *TemplateTest) TestHPACanSetCommonLabels() {
	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled":      "true",
			"commonLabels.environment": "production",
			"commonLabels.team":        "platform",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/hpa.yaml"})
	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(s.T(), output, &hpa)

	s.Require().Equal("production", hpa.Labels["environment"])
	s.Require().Equal("platform", hpa.Labels["team"])
}
