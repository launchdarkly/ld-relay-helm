package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	networkingv1 "k8s.io/api/networking/v1"
)

func (s *TemplateTest) TestIngressCanSetCommonLabels() {
	options := &helm.Options{
		SetValues: map[string]string{
			"ingress.enabled":          "true",
			"commonLabels.environment": "production",
			"commonLabels.team":        "platform",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/ingress.yaml"})
	var ingress networkingv1.Ingress
	helm.UnmarshalK8SYaml(s.T(), output, &ingress)

	s.Require().Equal("production", ingress.Labels["environment"])
	s.Require().Equal("platform", ingress.Labels["team"])
}
