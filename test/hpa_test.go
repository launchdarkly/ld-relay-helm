package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
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

func (s *TemplateTest) TestHPAAdditionalMetrics() {
	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled":                                           "true",
			"autoscaling.targetMemoryUtilizationPercentage":                 "0",
			"autoscaling.additionalMetrics[0].type":                         "External",
			"autoscaling.additionalMetrics[0].external.metric.name":         "datadogmetric@relay:relay-connections",
			"autoscaling.additionalMetrics[0].external.target.type":         "AverageValue",
			"autoscaling.additionalMetrics[0].external.target.averageValue": "2500",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/hpa.yaml"})
	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(s.T(), output, &hpa)

	s.Require().Len(hpa.Spec.Metrics, 2)

	s.Require().Equal(autoscalingv2.ResourceMetricSourceType, hpa.Spec.Metrics[0].Type)
	s.Require().Equal(corev1.ResourceCPU, hpa.Spec.Metrics[0].Resource.Name)

	s.Require().Equal(autoscalingv2.ExternalMetricSourceType, hpa.Spec.Metrics[1].Type)
	s.Require().NotNil(hpa.Spec.Metrics[1].External)
	s.Require().Equal("datadogmetric@relay:relay-connections", hpa.Spec.Metrics[1].External.Metric.Name)
	s.Require().Equal(autoscalingv2.AverageValueMetricType, hpa.Spec.Metrics[1].External.Target.Type)
	s.Require().Equal("2500", hpa.Spec.Metrics[1].External.Target.AverageValue.String())
}

func (s *TemplateTest) TestHPAAdditionalMetricsOnly() {
	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled":                                           "true",
			"autoscaling.targetCPUUtilizationPercentage":                    "0",
			"autoscaling.targetMemoryUtilizationPercentage":                 "0",
			"autoscaling.additionalMetrics[0].type":                         "External",
			"autoscaling.additionalMetrics[0].external.metric.name":         "datadogmetric@relay:relay-connections",
			"autoscaling.additionalMetrics[0].external.target.type":         "AverageValue",
			"autoscaling.additionalMetrics[0].external.target.averageValue": "2500",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/hpa.yaml"})
	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(s.T(), output, &hpa)

	s.Require().Len(hpa.Spec.Metrics, 1)
	s.Require().Equal(autoscalingv2.ExternalMetricSourceType, hpa.Spec.Metrics[0].Type)
}

// TestHPADefaultsToCPUWhenNoMetricsConfigured verifies that disabling every
// metric source renders 80% average CPU explicitly, rather than an empty
// `spec.metrics` that the API server silently defaults to the same thing.
func (s *TemplateTest) TestHPADefaultsToCPUWhenNoMetricsConfigured() {
	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled":                           "true",
			"autoscaling.targetCPUUtilizationPercentage":    "0",
			"autoscaling.targetMemoryUtilizationPercentage": "0",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/hpa.yaml"})
	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(s.T(), output, &hpa)

	s.Require().Len(hpa.Spec.Metrics, 1)
	s.Require().Equal(autoscalingv2.ResourceMetricSourceType, hpa.Spec.Metrics[0].Type)
	s.Require().Equal(corev1.ResourceCPU, hpa.Spec.Metrics[0].Resource.Name)
	s.Require().Equal(autoscalingv2.UtilizationMetricType, hpa.Spec.Metrics[0].Resource.Target.Type)
	s.Require().NotNil(hpa.Spec.Metrics[0].Resource.Target.AverageUtilization)
	s.Require().Equal(int32(80), *hpa.Spec.Metrics[0].Resource.Target.AverageUtilization)
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
