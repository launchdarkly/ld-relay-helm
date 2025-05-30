package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	policyv1 "k8s.io/api/policy/v1"
)

func (s *TemplateTest) TestCanEnablePDB() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.disruptionBudget.enabled": "true",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal(0, pdb.Spec.MinAvailable.IntValue())
}

func (s *TemplateTest) TestCanEnablePDBWithDeprecatedProperty() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.distruptionBudget.enabled": "true",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal(0, pdb.Spec.MinAvailable.IntValue())
}

func (s *TemplateTest) TestCanChangePDBMaxUnavailableNumber() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.disruptionBudget.enabled":        "true",
			"pod.disruptionBudget.maxUnavailable": "2",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal(2, pdb.Spec.MaxUnavailable.IntValue())
}

func (s *TemplateTest) TestCanChangePDBMaxUnavailableNumberWithDeprecatedProperty() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.distruptionBudget.enabled":        "true",
			"pod.distruptionBudget.maxUnavailable": "2",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal(2, pdb.Spec.MaxUnavailable.IntValue())
}

func (s *TemplateTest) TestCanChangePDBMaxUnavailablePercent() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.disruptionBudget.enabled":        "true",
			"pod.disruptionBudget.maxUnavailable": "50%",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal("50%", pdb.Spec.MaxUnavailable.StrVal)
}

func (s *TemplateTest) TestCanChangePDBMaxUnavailablePercentWithDeprecatedProperty() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.distruptionBudget.enabled":        "true",
			"pod.distruptionBudget.maxUnavailable": "50%",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal("50%", pdb.Spec.MaxUnavailable.StrVal)
}

func (s *TemplateTest) TestCanChangePDBMinAvailableNumber() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.disruptionBudget.enabled":      "true",
			"pod.disruptionBudget.minAvailable": "2",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal(2, pdb.Spec.MinAvailable.IntValue())
}

func (s *TemplateTest) TestCanChangePDBMinAvailableNumberWithDeprecatedProperty() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.distruptionBudget.enabled":      "true",
			"pod.distruptionBudget.minAvailable": "2",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal(2, pdb.Spec.MinAvailable.IntValue())
}

func (s *TemplateTest) TestCanChangePDBMinAvailablePercent() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.disruptionBudget.enabled":      "true",
			"pod.disruptionBudget.minAvailable": "50%",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal("50%", pdb.Spec.MinAvailable.StrVal)
}

func (s *TemplateTest) TestCanChangePDBMinAvailablePercentWithDeprecatedProperty() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.distruptionBudget.enabled":      "true",
			"pod.distruptionBudget.minAvailable": "50%",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal("50%", pdb.Spec.MinAvailable.StrVal)
}

func (s *TemplateTest) TestBothPDBValues() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.disruptionBudget.enabled":        "true",
			"pod.disruptionBudget.minAvailable":   "50%",
			"pod.disruptionBudget.maxUnavailable": "2",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal(2, pdb.Spec.MaxUnavailable.IntValue())
	s.Require().Nil(pdb.Spec.MinAvailable)
}

func (s *TemplateTest) TestBothPDBValuesWithDeprecatedProperty() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.distruptionBudget.enabled":        "true",
			"pod.distruptionBudget.minAvailable":   "50%",
			"pod.distruptionBudget.maxUnavailable": "2",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal(2, pdb.Spec.MaxUnavailable.IntValue())
	s.Require().Nil(pdb.Spec.MinAvailable)
}

func (s *TemplateTest) TestPDBCanSetCommonLabels() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.disruptionBudget.enabled": "true",
			"commonLabels.environment":     "production",
			"commonLabels.team":            "platform",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/poddisruptionbudget.yaml"})
	var pdb policyv1.PodDisruptionBudget
	helm.UnmarshalK8SYaml(s.T(), output, &pdb)

	s.Require().Equal("production", pdb.Labels["environment"])
	s.Require().Equal("platform", pdb.Labels["team"])
}
