package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
)

func (s *TemplateTest) TestHPANamespace() {
	s.assertTemplateNamespaces("HorizontalPodAutoscaler", "templates/hpa.yaml", map[string]string{
		"autoscaling.enabled": "true",
	})
}

func (s *TemplateTest) TestHPABehaviorNotSetByDefault() {
	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled": "true",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/hpa.yaml"})
	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(s.T(), output, &hpa)

	s.Require().Nil(hpa.Spec.Behavior)
}

func (s *TemplateTest) TestCanSetHPAScaleUpBehavior() {
	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled":                                    "true",
			"autoscaling.behavior.scaleUp.policies[0].type":          "Pods",
			"autoscaling.behavior.scaleUp.policies[0].value":         "6",
			"autoscaling.behavior.scaleUp.policies[0].periodSeconds": "120",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/hpa.yaml"})
	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(s.T(), output, &hpa)

	s.Require().NotNil(hpa.Spec.Behavior)
	s.Require().NotNil(hpa.Spec.Behavior.ScaleUp)
	s.Require().Len(hpa.Spec.Behavior.ScaleUp.Policies, 1)
	s.Require().Equal(autoscalingv2.PodsScalingPolicy, hpa.Spec.Behavior.ScaleUp.Policies[0].Type)
	s.Require().Equal(int32(6), hpa.Spec.Behavior.ScaleUp.Policies[0].Value)
	s.Require().Equal(int32(120), hpa.Spec.Behavior.ScaleUp.Policies[0].PeriodSeconds)
}

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
