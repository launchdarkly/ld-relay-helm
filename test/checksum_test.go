package test

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	appsv1 "k8s.io/api/apps/v1"
)

// The checksum/config annotation rolls the pods when the ConfigMap data changes.
// It must not react to the chart metadata. The ConfigMap labels hold the chart
// version, so a checksum over the whole manifest restarts every pod on each
// chart release. These tests pin both halves of that behavior.

// renderConfigChecksum renders the deployment and returns the checksum/config
// annotation on the pod template.
func (s *TemplateTest) renderConfigChecksum(chartPath string, setValues map[string]string) string {
	options := &helm.Options{
		SetValues:      setValues,
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, chartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	checksum := deployment.Spec.Template.Annotations["checksum/config"]
	s.Require().NotEmpty(checksum, "the pod template must carry a checksum/config annotation")

	return checksum
}

// chartCopyWithVersion copies the chart to a temporary directory and replaces
// the version and appVersion. The copy lets a test render the same values
// against a different chart version.
func (s *TemplateTest) chartCopyWithVersion(version, appVersion string) string {
	dest := filepath.Join(s.T().TempDir(), "chart")
	s.Require().NoError(os.MkdirAll(dest, 0o755))

	templates := filepath.Join(s.ChartPath, "templates")
	s.Require().NoError(os.CopyFS(filepath.Join(dest, "templates"), os.DirFS(templates)))

	values, err := os.ReadFile(filepath.Join(s.ChartPath, "values.yaml"))
	s.Require().NoError(err)
	s.Require().NoError(os.WriteFile(filepath.Join(dest, "values.yaml"), values, 0o644))

	metadata, err := os.ReadFile(filepath.Join(s.ChartPath, "Chart.yaml"))
	s.Require().NoError(err)
	metadata = regexp.MustCompile(`(?m)^version:.*$`).ReplaceAll(metadata, []byte("version: "+version))
	metadata = regexp.MustCompile(`(?m)^appVersion:.*$`).ReplaceAll(metadata, []byte(`appVersion: "`+appVersion+`"`))
	s.Require().NoError(os.WriteFile(filepath.Join(dest, "Chart.yaml"), metadata, 0o644))

	return dest
}

func (s *TemplateTest) TestChecksumConfigIgnoresChartVersion() {
	setValues := map[string]string{
		"relay.environment.LD_ENV_test": "sdk-key",
	}

	current := s.renderConfigChecksum(s.ChartPath, setValues)
	bumped := s.renderConfigChecksum(s.chartCopyWithVersion("99.99.99", "99.99.99"), setValues)

	s.Require().Equal(current, bumped,
		"a chart version bump must not change checksum/config, or every helm upgrade restarts the pods")
}

func (s *TemplateTest) TestChecksumConfigIgnoresConfigMapLabels() {
	setValues := map[string]string{
		"relay.environment.LD_ENV_test": "sdk-key",
	}

	labelled := map[string]string{
		"relay.environment.LD_ENV_test": "sdk-key",
		"commonLabels.team":             "sdk",
	}

	s.Require().Equal(s.renderConfigChecksum(s.ChartPath, setValues), s.renderConfigChecksum(s.ChartPath, labelled),
		"checksum/config must cover the ConfigMap data only, not its metadata")
}

func (s *TemplateTest) TestChecksumConfigChangesWithConfigData() {
	first := s.renderConfigChecksum(s.ChartPath, map[string]string{
		"relay.environment.LD_ENV_test": "sdk-key-one",
	})

	second := s.renderConfigChecksum(s.ChartPath, map[string]string{
		"relay.environment.LD_ENV_test": "sdk-key-two",
	})

	s.Require().NotEqual(first, second, "a config change must change checksum/config so the pods roll")
}

func (s *TemplateTest) TestChecksumConfigChangesWhenEnvironmentIsAdded() {
	first := s.renderConfigChecksum(s.ChartPath, map[string]string{
		"relay.environment.LD_ENV_test": "sdk-key",
	})

	second := s.renderConfigChecksum(s.ChartPath, map[string]string{
		"relay.environment.LD_ENV_test":  "sdk-key",
		"relay.environment.LD_ENV_other": "other-sdk-key",
	})

	s.Require().NotEqual(first, second, "an added environment must change checksum/config so the pods roll")
}
