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

func (s *TemplateTest) TestCanSetEnvironmentVariablesUsingEnvironmentVariables() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.environmentVariables[0].name":                         "RUNNING_HOST",
			"relay.environmentVariables[0].valueFrom.fieldRef.fieldPath": "status.hostIP",
			"relay.environmentVariables[1].name":                         "ENVIRONMENT",
			"relay.environmentVariables[1].value":                        "production",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("RUNNING_HOST", deployment.Spec.Template.Spec.Containers[0].Env[0].Name)
	s.Require().Equal("status.hostIP", deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.FieldRef.FieldPath)

	s.Require().Equal("ENVIRONMENT", deployment.Spec.Template.Spec.Containers[0].Env[1].Name)
	s.Require().Equal("production", deployment.Spec.Template.Spec.Containers[0].Env[1].Value)
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

func (s *TemplateTest) TestCanSetEnvironmentVariablesFromEnvFromSecrets() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.secretEnvironmentVariables": "true",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("ld-relay-test-secret-environment-variables", deployment.Spec.Template.Spec.Containers[0].EnvFrom[0].SecretRef.Name)
}

func (s *TemplateTest) TestCanMountSecretsAsVolumes() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.secrets[0].volumePath": "my-secret-path",
			"relay.secrets[0].secretKey":  "ld-relay",
			"relay.secrets[0].secretName": "sdk-key",
			"relay.secrets[1].volumePath": "my-second-secret-path",
			"relay.secrets[1].secretKey":  "id",
			"relay.secrets[1].secretName": "application-info",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal(0, len(deployment.Spec.Template.Spec.Containers[0].Env))

	s.Require().Len(deployment.Spec.Template.Spec.Containers[0].VolumeMounts, 1)
	s.Require().Equal("ld-relay-projected-volume", deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name)
	s.Require().Equal("/mnt/secrets/", deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath)

	s.Require().Len(deployment.Spec.Template.Spec.Volumes, 2)

	volume := deployment.Spec.Template.Spec.Volumes[1]
	s.Require().Equal("ld-relay-projected-volume", volume.Name)

	s.Require().Len(volume.Projected.Sources, 2)
	s.Require().Equal("sdk-key", volume.Projected.Sources[0].Secret.Name)
	s.Require().Equal("ld-relay", volume.Projected.Sources[0].Secret.Items[0].Key)
	s.Require().Equal("my-secret-path", volume.Projected.Sources[0].Secret.Items[0].Path)

	s.Require().Equal("application-info", volume.Projected.Sources[1].Secret.Name)
	s.Require().Equal("id", volume.Projected.Sources[1].Secret.Items[0].Key)
	s.Require().Equal("my-second-secret-path", volume.Projected.Sources[1].Secret.Items[0].Path)
}

func (s *TemplateTest) TestCanLoadConfigFromVolume() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.volume.config": "my-new-config.config",
			"relay.volume.definition.persistentVolumeClaim.claimName": "ld-relay-pvc",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("ld-relay-config", deployment.Spec.Template.Spec.Volumes[0].Name)

	s.Require().Equal("ld-relay-volume", deployment.Spec.Template.Spec.Volumes[1].Name)
	s.Require().Equal("ld-relay-pvc", deployment.Spec.Template.Spec.Volumes[1].PersistentVolumeClaim.ClaimName)

	s.Require().Equal("ld-relay-volume", deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name)
	s.Require().Equal("/mnt/volume", deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath)

	expectedCommand := []string{
		"/usr/bin/ldr",
		"--config",
		"/mnt/volume/my-new-config.config",
		"--allow-missing-file",
		"--from-env",
	}
	s.Require().Equal(expectedCommand, deployment.Spec.Template.Spec.Containers[0].Command)
}

func (s *TemplateTest) TestCanEnableOfflineMode() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.volume.offline": "relay-file.tar.gz",
			"relay.volume.definition.persistentVolumeClaim.claimName": "ld-relay-pvc",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("ld-relay-config", deployment.Spec.Template.Spec.Volumes[0].Name)

	s.Require().Equal("ld-relay-volume", deployment.Spec.Template.Spec.Volumes[1].Name)
	s.Require().Equal("ld-relay-pvc", deployment.Spec.Template.Spec.Volumes[1].PersistentVolumeClaim.ClaimName)

	s.Require().Equal("ld-relay-volume", deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name)
	s.Require().Equal("/mnt/volume", deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath)

	s.Require().Equal("FILE_DATA_SOURCE", deployment.Spec.Template.Spec.Containers[0].Env[0].Name)
	s.Require().Equal("/mnt/volume/relay-file.tar.gz", deployment.Spec.Template.Spec.Containers[0].Env[0].Value)
}

func (s *TemplateTest) TestCanSetPodAnnotations() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.annotations.first-annotation":  "example-value-one",
			"pod.annotations.second-annotation": "example-value-two",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("example-value-one", deployment.Spec.Template.Annotations["first-annotation"])
	s.Require().Equal("example-value-two", deployment.Spec.Template.Annotations["second-annotation"])
}

func (s *TemplateTest) TestCanSetDeprecatedPodAnnotations() {
	options := &helm.Options{
		SetValues: map[string]string{
			"podAnnotations.first-annotation":  "example-value-one",
			"podAnnotations.second-annotation": "example-value-two",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("example-value-one", deployment.Spec.Template.Annotations["first-annotation"])
	s.Require().Equal("example-value-two", deployment.Spec.Template.Annotations["second-annotation"])
}

func (s *TemplateTest) TestPodAnnotationsTakesPrecendenceOverDeprecatedOption() {
	options := &helm.Options{
		SetValues: map[string]string{
			"podAnnotations.testing-annotation":  "legacy",
			"pod.annotations.testing-annotation": "new",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("new", deployment.Spec.Template.Annotations["testing-annotation"])
}

func (s *TemplateTest) TestCanSetPodLabels() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.labels.first-label":  "example-value-one",
			"pod.labels.second-label": "example-value-two",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("example-value-one", deployment.Spec.Template.Labels["first-label"])
	s.Require().Equal("example-value-two", deployment.Spec.Template.Labels["second-label"])
}

func (s *TemplateTest) TestCanSetPodSecurityContext() {
	options := &helm.Options{
		SetValues: map[string]string{
			"pod.securityContext.runAsUser":  "1000",
			"pod.securityContext.runAsGroup": "2000",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal(int64(1000), *deployment.Spec.Template.Spec.SecurityContext.RunAsUser)
	s.Require().Equal(int64(2000), *deployment.Spec.Template.Spec.SecurityContext.RunAsGroup)
}

func (s *TemplateTest) TestProbesDefaultToSaneValues() {
	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}
	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	var none corev1.URIScheme
	none = ""
	s.Require().Equal("/status", *&deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Path)
	s.Require().Equal("api", deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Port.String())
	s.Require().Equal(none, *&deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Scheme)

	s.Require().Equal("/status", *&deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Path)
	s.Require().Equal("api", deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Port.String())
	s.Require().Equal(none, *&deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Scheme)
}

func (s *TemplateTest) TestCanAffectHttpGetProbes() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.livenessProbe.httpGet.path":   "/liveness",
			"relay.livenessProbe.httpGet.port":   "8000",
			"relay.livenessProbe.httpGet.scheme": "HTTPS",

			"relay.readinessProbe.httpGet.path":   "/readiness",
			"relay.readinessProbe.httpGet.port":   "9000",
			"relay.readinessProbe.httpGet.scheme": "HTTPS",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal("/liveness", *&deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Path)
	s.Require().Equal(int(8000), deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Port.IntValue())
	s.Require().Equal(corev1.URISchemeHTTPS, *&deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Scheme)

	s.Require().Equal("/readiness", *&deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Path)
	s.Require().Equal(int(9000), deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Port.IntValue())
	s.Require().Equal(corev1.URISchemeHTTPS, *&deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Scheme)

}

func (s *TemplateTest) TestCanDisableProbes() {
	options := &helm.Options{
		SetValues: map[string]string{
			"relay.livenessProbe":  "null",
			"relay.readinessProbe": "null",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Nil(deployment.Spec.Template.Spec.Containers[0].LivenessProbe)
	s.Require().Nil(deployment.Spec.Template.Spec.Containers[0].ReadinessProbe)
}

func (s *TemplateTest) TestCanSetDeprecatedPodSecurityContext() {
	options := &helm.Options{
		SetValues: map[string]string{
			"podSecurityContext.runAsUser":  "1000",
			"podSecurityContext.runAsGroup": "2000",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal(int64(1000), *deployment.Spec.Template.Spec.SecurityContext.RunAsUser)
	s.Require().Equal(int64(2000), *deployment.Spec.Template.Spec.SecurityContext.RunAsGroup)
}

func (s *TemplateTest) TestPodSecurityContextTakesPrecendenceOverDeprecatedOption() {
	options := &helm.Options{
		SetValues: map[string]string{
			"podSecurityContext.runAsUser":  "1000",
			"pod.securityContext.runAsUser": "2000",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.Namespace),
	}

	output := helm.RenderTemplate(s.T(), options, s.ChartPath, s.Release, []string{"templates/deployment.yaml"})
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	s.Require().Equal(int64(2000), *deployment.Spec.Template.Spec.SecurityContext.RunAsUser)
}
