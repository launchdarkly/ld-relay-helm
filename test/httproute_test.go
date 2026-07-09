package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func (s *TemplateTest) TestHTTPRouteNamespace() {
	s.assertTemplateNamespaces("HTTPRoute", "templates/httproute.yaml", map[string]string{
		"httpRoute.enabled":                        "true",
		"httpRoute.hostnames[0]":                   "ld-relay.local",
		"httpRoute.rules[0].matches[0].path.type":  "PathPrefix",
		"httpRoute.rules[0].matches[0].path.value": "/",
		"httpRoute.rules[0].port":                  "8030",
	})
}

type httpRouteLabels struct {
	Metadata struct {
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

func (s *TemplateTest) TestHTTPRouteCanSetCommonLabels() {
	options := &helm.Options{
		SetValues: map[string]string{
			"httpRoute.enabled":                        "true",
			"httpRoute.rules[0].matches[0].path.type":  "PathPrefix",
			"httpRoute.rules[0].matches[0].path.value": "/",
			"httpRoute.rules[0].port":                  "8030",
			"commonLabels.environment":                 "production",
			"commonLabels.team":                        "platform",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/httproute.yaml"})
	var route httpRouteLabels
	helm.UnmarshalK8SYaml(s.T(), output, &route)

	s.Require().Equal("production", route.Metadata.Labels["environment"])
	s.Require().Equal("platform", route.Metadata.Labels["team"])
}
