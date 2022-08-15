package test

import (
	"path/filepath"
	"strings"
	"testing"

	"ld-relay-helm/test/golden"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestGoldenDefaultsTemplate(t *testing.T) {
	t.Parallel()

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)
	templateNames := []string{"config", "deployment", "service", "serviceaccount"}

	for _, name := range templateNames {
		suite.Run(t, &golden.TemplateGoldenTest{
			ChartPath:      chartPath,
			Release:        "ld-relay-test",
			Namespace:      "ld-relay-" + strings.ToLower(random.UniqueId()),
			GoldenFileName: name,
			Templates:      []string{"templates/" + name + ".yaml"},
		})
	}
}

func TestGoldenIngressWithBaseConfiguration(t *testing.T) {
	t.Parallel()

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	suite.Run(t, &golden.TemplateGoldenTest{
		ChartPath:      chartPath,
		Release:        "ld-relay-test",
		Namespace:      "ld-relay-" + strings.ToLower(random.UniqueId()),
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
