package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	corev1 "k8s.io/api/core/v1"
)

func (s *TemplateTest) TestServiceSupportsMultiplePorts() {
	options := &helm.Options{
		SetValues: map[string]string{
			"service.ports[0].port":       "80",
			"service.ports[0].protocol":   "TCP",
			"service.ports[0].name":       "api",
			"service.ports[0].targetPort": "8080",
			"service.ports[1].port":       "88",
			"service.ports[1].protocol":   "TCP",
			"service.ports[1].name":       "prometheus",
			"service.ports[1].targetPort": "8088",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/service.yaml"})
	var service corev1.Service
	helm.UnmarshalK8SYaml(s.T(), output, &service)

	s.Require().Equal(2, len(service.Spec.Ports))

	s.Require().Equal(int32(80), service.Spec.Ports[0].Port)
	s.Require().Equal(8080, service.Spec.Ports[0].TargetPort.IntValue())
	s.Require().Equal(corev1.Protocol("TCP"), service.Spec.Ports[0].Protocol)
	s.Require().Equal("api", service.Spec.Ports[0].Name)

	s.Require().Equal(int32(88), service.Spec.Ports[1].Port)
	s.Require().Equal(8088, service.Spec.Ports[1].TargetPort.IntValue())
	s.Require().Equal(corev1.Protocol("TCP"), service.Spec.Ports[1].Protocol)
	s.Require().Equal("prometheus", service.Spec.Ports[1].Name)
}
