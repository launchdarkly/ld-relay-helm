package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (s *TemplateTest) TestCanControlContainerPorts() {
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

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal(int32(8080), deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort)
	s.Require().Equal(corev1.Protocol("TCP"), deployment.Spec.Template.Spec.Containers[0].Ports[0].Protocol)
	s.Require().Equal("api", deployment.Spec.Template.Spec.Containers[0].Ports[0].Name)

	s.Require().Equal(int32(8088), deployment.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort)
	s.Require().Equal(corev1.Protocol("TCP"), deployment.Spec.Template.Spec.Containers[0].Ports[1].Protocol)
	s.Require().Equal("prometheus", deployment.Spec.Template.Spec.Containers[0].Ports[1].Name)
}

func (s *TemplateTest) TestCanSetEnvironmentVariablesFromSecrets() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.secrets[0].envName":    "LD_ENV_Production",
			"relay.secrets[0].secretKey":  "ld-relay",
			"relay.secrets[0].secretName": "sdk-key",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("LD_ENV_Production", deployment.Spec.Template.Spec.Containers[0].Env[0].Name)
	s.Require().Equal("ld-relay", deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.Key)
	s.Require().Equal("sdk-key", deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.Name)
}

func (s *TemplateTest) TestCanEnableOfflineMode() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.offline.enabled":                                "true",
			"relay.offline.filename":                               "relay-file.tar.gz",
			"relay.offline.volume.persistentVolumeClaim.claimName": "ld-relay-offline-pvc",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("ld-relay-config", deployment.Spec.Template.Spec.Volumes[0].Name)

	s.Require().Equal("ld-relay-offline", deployment.Spec.Template.Spec.Volumes[1].Name)
	s.Require().Equal("ld-relay-offline-pvc", deployment.Spec.Template.Spec.Volumes[1].PersistentVolumeClaim.ClaimName)

	s.Require().Equal("ld-relay-offline", deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name)
	s.Require().Equal("/offline/", deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath)

	s.Require().Equal("FILE_DATA_SOURCE", deployment.Spec.Template.Spec.Containers[0].Env[0].Name)
	s.Require().Equal("/offline/relay-file.tar.gz", deployment.Spec.Template.Spec.Containers[0].Env[0].Value)
}
