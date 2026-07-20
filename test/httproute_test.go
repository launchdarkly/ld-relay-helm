package test

import (
	"os/exec"
	"path/filepath"

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

type httpRouteSpec struct {
	Metadata struct {
		Name        string            `json:"name"`
		Annotations map[string]string `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Hostnames []string `json:"hostnames"`
		Rules     []struct {
			BackendRefs []struct {
				Name   string `json:"name"`
				Port   int    `json:"port"`
				Weight int    `json:"weight"`
			} `json:"backendRefs"`
		} `json:"rules"`
	} `json:"spec"`
}

// TestHTTPRouteDefaultBackendUsesChartService verifies the `port` convenience:
// when a rule gives only a port (no explicit backendRefs), the backend is
// synthesized to point at this chart's own Service. The Service name is the
// computed fullname, which is also the HTTPRoute's own name, so the two must
// match.
func (s *TemplateTest) TestHTTPRouteDefaultBackendUsesChartService() {
	options := &helm.Options{
		SetValues: map[string]string{
			"httpRoute.enabled":                        "true",
			"httpRoute.rules[0].matches[0].path.type":  "PathPrefix",
			"httpRoute.rules[0].matches[0].path.value": "/",
			"httpRoute.rules[0].port":                  "8030",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/httproute.yaml"})
	var route httpRouteSpec
	helm.UnmarshalK8SYaml(s.T(), output, &route)

	s.Require().Len(route.Spec.Rules, 1)
	s.Require().Len(route.Spec.Rules[0].BackendRefs, 1)
	backend := route.Spec.Rules[0].BackendRefs[0]
	s.Require().Equal(route.Metadata.Name, backend.Name, "default backend must point at this chart's own Service")
	s.Require().Equal(8030, backend.Port)
}

// TestHTTPRouteAnnotations verifies httpRoute.annotations render onto the
// resource metadata.
func (s *TemplateTest) TestHTTPRouteAnnotations() {
	options := &helm.Options{
		SetValues: map[string]string{
			"httpRoute.enabled":                        "true",
			"httpRoute.annotations.example\\.com/team": "platform",
			"httpRoute.rules[0].matches[0].path.type":  "PathPrefix",
			"httpRoute.rules[0].matches[0].path.value": "/",
			"httpRoute.rules[0].port":                  "8030",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/httproute.yaml"})
	var route httpRouteSpec
	helm.UnmarshalK8SYaml(s.T(), output, &route)

	s.Require().Equal("platform", route.Metadata.Annotations["example.com/team"])
}

// TestHTTPRouteHostnamesQuoting verifies hostnames render as valid YAML,
// including a wildcard hostname that requires quoting (a bare `- *.foo` would be
// an invalid YAML alias). Successful unmarshal plus the preserved value proves
// the toYaml rendering quotes where needed.
func (s *TemplateTest) TestHTTPRouteHostnamesQuoting() {
	options := &helm.Options{
		SetValues: map[string]string{
			"httpRoute.enabled":       "true",
			"httpRoute.hostnames[0]":  "*.ld-relay.local",
			"httpRoute.hostnames[1]":  "ld-relay.local",
			"httpRoute.rules[0].port": "8030",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/httproute.yaml"})
	var route httpRouteSpec
	helm.UnmarshalK8SYaml(s.T(), output, &route)

	s.Require().Equal([]string{"*.ld-relay.local", "ld-relay.local"}, route.Spec.Hostnames)
}

// TestHTTPRoutePathlessMatch verifies that a match without a `path` (a valid
// Gateway API match that constrains only headers/method/queryParams) renders
// instead of crashing the template with a nil-pointer dereference.
func (s *TemplateTest) TestHTTPRoutePathlessMatch() {
	options := &helm.Options{
		SetValues: map[string]string{
			"httpRoute.enabled":                              "true",
			"httpRoute.parentRefs[0].name":                   "my-gateway",
			"httpRoute.rules[0].matches[0].headers[0].name":  "X-Env",
			"httpRoute.rules[0].matches[0].headers[0].value": "prod",
			"httpRoute.rules[0].port":                        "8030",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output, err := helm.RenderTemplateE(s.T(), options, s.ChartPath, s.Release,
		[]string{"templates/httproute.yaml"})

	s.Require().NoError(err, "a header-only match (no path) must not crash the template")
	s.Require().Contains(output, "kind: HTTPRoute")
}

// TestHTTPRouteBackendRefsOverride verifies that an explicit backendRefs takes
// precedence over the `port` convenience default (which points at the chart
// Service), enabling weighting / mirroring / external backends.
func (s *TemplateTest) TestHTTPRouteBackendRefsOverride() {
	options := &helm.Options{
		SetValues: map[string]string{
			"httpRoute.enabled":                        "true",
			"httpRoute.rules[0].matches[0].path.type":  "PathPrefix",
			"httpRoute.rules[0].matches[0].path.value": "/",
			"httpRoute.rules[0].backendRefs[0].name":   "external-svc",
			"httpRoute.rules[0].backendRefs[0].port":   "9000",
			"httpRoute.rules[0].backendRefs[0].weight": "100",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/httproute.yaml"})
	var route httpRouteSpec
	helm.UnmarshalK8SYaml(s.T(), output, &route)

	s.Require().Len(route.Spec.Rules, 1)
	s.Require().Len(route.Spec.Rules[0].BackendRefs, 1)
	backend := route.Spec.Rules[0].BackendRefs[0]
	s.Require().Equal("external-svc", backend.Name, "explicit backendRefs must not be overwritten by the chart Service default")
	s.Require().Equal(9000, backend.Port)
	s.Require().Equal(100, backend.Weight)
}

// TestHTTPRouteRuleRequiresBackend verifies the guard that fails the render when
// a rule specifies neither a `port` nor explicit `backendRefs`, surfacing the
// mistake at template time instead of rendering an invalid null backend port.
func (s *TemplateTest) TestHTTPRouteRuleRequiresBackend() {
	options := &helm.Options{
		SetValues: map[string]string{
			"httpRoute.enabled":                        "true",
			"httpRoute.rules[0].matches[0].path.type":  "PathPrefix",
			"httpRoute.rules[0].matches[0].path.value": "/",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	_, err := helm.RenderTemplateE(s.T(), options, s.ChartPath, s.Release,
		[]string{"templates/httproute.yaml"})

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "needs a `port`")
}

// TestHTTPRouteParentRefsPassThrough verifies parentRefs render the full
// ParentReference schema -- including `port` and `sectionName`, which the
// previous hand-rendered template dropped.
func (s *TemplateTest) TestHTTPRouteParentRefsPassThrough() {
	options := &helm.Options{
		SetValues: map[string]string{
			"httpRoute.enabled":                        "true",
			"httpRoute.parentRefs[0].name":             "my-gateway",
			"httpRoute.parentRefs[0].namespace":        "gateway-infra",
			"httpRoute.parentRefs[0].sectionName":      "https",
			"httpRoute.parentRefs[0].port":             "443",
			"httpRoute.rules[0].matches[0].path.type":  "PathPrefix",
			"httpRoute.rules[0].matches[0].path.value": "/",
			"httpRoute.rules[0].port":                  "8030",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/httproute.yaml"})

	s.Require().Contains(output, "sectionName: https")
	s.Require().Contains(output, "port: 443")
}

// TestHTTPRoutePortNotInjectable verifies that a rule's `port` is coerced to an
// integer, so a string port carrying a YAML-injection payload (newline + a fake
// sibling key) cannot inject keys into the rendered backendRef.
func (s *TemplateTest) TestHTTPRoutePortNotInjectable() {
	valuesFile, err := filepath.Abs("testdata/httproute_port_injection.yaml")
	s.Require().NoError(err)

	options := &helm.Options{
		ValuesFiles:    []string{valuesFile},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/httproute.yaml"})

	s.Require().NotContains(output, "injected")
	s.Require().NotContains(output, "pwned")
	var route httpRouteSpec
	helm.UnmarshalK8SYaml(s.T(), output, &route)
	s.Require().Len(route.Spec.Rules, 1)
	s.Require().Len(route.Spec.Rules[0].BackendRefs, 1)
}

// TestHTTPRouteNotesHeaderOnlyMatch renders NOTES.txt (via `helm install
// --dry-run`, since `helm template` omits NOTES) for a header-only match. This
// locks in the NOTES `.path`/`.value` guard against regression: a match with no
// path must not nil-panic, must reach the httproute guidance, and must not fall
// through to the ClusterIP port-forward branch.
func (s *TemplateTest) TestHTTPRouteNotesHeaderOnlyMatch() {
	cmd := exec.Command("helm", "install", s.Release, s.ChartPath,
		"--dry-run=client", "--namespace", s.Namespace,
		"--set", "httpRoute.enabled=true",
		"--set", "httpRoute.rules[0].matches[0].headers[0].name=X-Env",
		"--set", "httpRoute.rules[0].matches[0].headers[0].value=prod",
		"--set", "httpRoute.rules[0].port=8030",
	)
	out, err := cmd.CombinedOutput()
	s.Require().NoError(err, string(out))

	notes := string(out)
	// Assert on text unique to the HTTPRoute NOTES branch (not just "httproute",
	// which also appears in the rendered manifest dump), and confirm we did not
	// fall through to the ClusterIP port-forward branch.
	s.Require().Contains(notes, "configured on your Gateway")
	s.Require().NotContains(notes, "port-forward")
}
