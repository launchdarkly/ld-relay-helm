package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const goldenNamespace = "ld-relay-test-ns"

func TestGoldenDefaultsTemplate(t *testing.T) {
	t.Parallel()

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"config", "deployment", "service", "serviceaccount"}

	for _, name := range templateNames {
		suite.Run(t, &TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "ld-relay-test",
			Namespace:      goldenNamespace,
			GoldenFileName: name,
			Templates:      []string{"templates/" + name + ".yaml"},
		})
	}
}

func TestGoldenIngressWithBaseConfiguration(t *testing.T) {
	t.Parallel()

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	suite.Run(t, &TemplateGoldenTest{
		ChartPath:      chartPath,
		Release:        "ld-relay-test",
		Namespace:      goldenNamespace,
		GoldenFileName: "ingress",
		Templates:      []string{"templates/ingress.yaml"},
		SetValues: map[string]string{
			"ingress.enabled":       "true",
			"ingress.hosts[0].host": "ld-relay.local",

			"ingress.hosts[0].paths[0].path":     "/api",
			"ingress.hosts[0].paths[0].pathType": "Prefix",
			"ingress.hosts[0].paths[0].port":     "8030",

			"ingress.hosts[0].paths[1].path":     "/prometheus",
			"ingress.hosts[0].paths[1].pathType": "Prefix",
			"ingress.hosts[0].paths[1].port":     "8031",
		},
	})
}

func TestGoldenHTTPRouteWithBaseConfiguration(t *testing.T) {
	t.Parallel()

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	suite.Run(t, &TemplateGoldenTest{
		ChartPath:      chartPath,
		Release:        "ld-relay-test",
		Namespace:      goldenNamespace,
		GoldenFileName: "httproute",
		Templates:      []string{"templates/httproute.yaml"},
		SetValues: map[string]string{
			"httpRoute.enabled":                 "true",
			"httpRoute.parentRefs[0].name":      "my-gateway",
			"httpRoute.parentRefs[0].namespace": "gateway-infra",
			"httpRoute.hostnames[0]":            "ld-relay.local",

			"httpRoute.rules[0].matches[0].path.type":  "PathPrefix",
			"httpRoute.rules[0].matches[0].path.value": "/api",
			"httpRoute.rules[0].port":                  "8030",

			"httpRoute.rules[1].matches[0].path.type":  "PathPrefix",
			"httpRoute.rules[1].matches[0].path.value": "/prometheus",
			"httpRoute.rules[1].port":                  "8031",
		},
	})
}
