package test

import (
	"maps"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

const (
	testHelmNamespace   = "test-helm-namespace"
	testValuesNamespace = "test-values-namespace"
)

type renderedResource struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Namespace string `json:"namespace"`
	} `json:"metadata"`
}

func (s *TemplateTest) assertTemplateNamespaces(kind string, template string, setValues map[string]string) {
	testCases := []struct {
		Name              string
		HelmNamespace     string
		ValueNamespace    string
		ExpectedNamespace string
	}{
		{
			Name:              "DefaultNamespace",
			ExpectedNamespace: "default",
		},
		{
			Name:              "HelmNamespace",
			HelmNamespace:     testHelmNamespace,
			ExpectedNamespace: testHelmNamespace,
		},
		{
			Name:              "NamespaceOverride",
			HelmNamespace:     testHelmNamespace,
			ValueNamespace:    testValuesNamespace,
			ExpectedNamespace: testValuesNamespace,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		s.Run(testCase.Name, func() {
			options := &helm.Options{
				SetValues: namespaceTestSetValues(setValues, testCase.ValueNamespace),
			}

			if testCase.HelmNamespace != "" {
				options.KubectlOptions = k8s.NewKubectlOptions("", "", testCase.HelmNamespace)
			}

			output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{template})

			var rendered renderedResource
			helm.UnmarshalK8SYaml(s.T(), output, &rendered)

			s.Require().Equal(kind, rendered.Kind)
			s.Require().Equal(testCase.ExpectedNamespace, rendered.Metadata.Namespace)
		})
	}
}

func namespaceTestSetValues(setValues map[string]string, valueNamespace string) map[string]string {
	allValues := map[string]string{}

	maps.Copy(allValues, setValues)

	allValues["namespaceOverride"] = valueNamespace

	return allValues
}
