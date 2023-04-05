# LaunchDarkly Relay Proxy Helm Chart

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/launchdarkly/ld-relay-helm/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/launchdarkly/ld-relay-helm/tree/main)

A Helm chart to ease deployment of the [LaunchDarkly Relay Proxy](https://github.com/launchdarkly/ld-relay) to Kubernetes (k8s).

Basic [installation](#installation) and [configuration](#configuration-options) information is below.

To learn more, read the [Getting started](./docs/getting-started.md) guide. For additional examples, read:

* [Offline mode](./docs/examples/offline-mode.md)
* [Automatic configuration](./docs/examples/automatic-configuration.md)

## LaunchDarkly overview

[LaunchDarkly](https://www.launchdarkly.com) is a feature management platform that serves over 100 billion feature flags daily to help teams build better software, faster. [Get started](https://docs.launchdarkly.com/docs/getting-started) using LaunchDarkly today!

[![Twitter Follow](https://img.shields.io/twitter/follow/launchdarkly.svg?style=social&label=Follow&maxAge=2592000)](https://twitter.com/intent/follow?screen_name=launchdarkly)

## Installation

The default configuration is insufficient to have a working instance of the Relay Proxy running. You must minimally provide an environment for the Relay Proxy to connect to, using your LaunchDarkly SDK key for that environment.

To deploy the Relay Proxy to the Kubernetes cluster using the default configuration and an environment:

```shell
helm repo add launchdarkly-ld-relay https://launchdarkly.github.io/ld-relay-helm
helm install ld-relay --set relay.environment.LD_ENV_YourEnvironment=your-sdk-key launchdarkly-ld-relay/ld-relay
```

For additional configuration, use the [Configuration options](#configuration-options) below.

## Configuration options

To customize this Helm chart, override the configuration options defined in the [values file](https://github.com/launchdarkly/ld-relay-helm/blob/main/values.yaml). The values file contains detailed documentation on each option.

Here's a summary of the available configuration options:


| Key                                           | Type    | Default                                                      | Description                                                                      |
|-----------------------------------------------|---------|--------------------------------------------------------------|----------------------------------------------------------------------------------|
| relay.environment                             | object  | `{}`                                                         | Defines container environment variables to configure the Relay Proxy instance    |
| relay.secrets                                 | array   | `[]`                                                         | Defines container environment variables or volumes populated from k8s secrets    |
| relay.volume                                  | object  | `{}`                                                         | Enables offline mode or references an existing config file from a defined volume |
| replicaCount                                  | integer | `1`                                                          | Number of replicas of the relay pod                                              |
| image.repository                              | string  | `launchdarkly/ld-relay`                                      | ld-relay image repository                                                        |
| image.pullPolicy                              | string  | `IfNotPresent`                                               | ld-relay image pull policy                                                       |
| image.tag                                     | string  | `""`                                                         | Overrides the image tag whose default is the chart appVersion                    |
| imagePullSecrets                              | array   | `[]`                                                         | Specifies docker registry secret names as an array                               |
| nameOverride                                  | string  | `""`                                                         | Partially overrides the fullname template with a string (includes release name)  |
| fullnameOverride                              | string  | `""`                                                         | Fully overrides the fullname template with a string                              |
| serviceAccount.create                         | bool    | `true`                                                       | Specifies whether a service account should be created                            |
| serviceAccount.annotations                    | object  | `{}`                                                         | Annotations to add to the service account                                        |
| serviceAccount.name                           | string  | `""`                                                         | The name of the service account                                                  |
| podAnnotations                                | object  | `{}`                                                         | Pod annotations (deprecated: use pod.annotations instead)                        |
| podSecurityContext                            | object  | `{}`                                                         | Pod security context (deprecated: use pod.securityContext instead)               |
| pod.annotations                               | object  | `{}`                                                         | Pod annotations                                                                  |
| pod.labels                                    | object  | `{}`                                                         | Pod labels                                                                       |
| pod.securityContext                           | object  | `{}`                                                         | Pod security context                                                             |
| securityContext                               | object  | `{}`                                                         | Container security context                                                       |
| service.type                                  | string  | `ClusterIP`                                                  | Kubernetes service type                                                          |
| service.annotations                           | object  | `{}`                                                         | Annotations to add to the service                                                |
| service.ports                                 | array   | `[{port: 8030, targetPort: 8030, protocol: TCP, name: api}]` | Service port mapping. Must include one port named `api`.                         |
| ingress.enabled                               | bool    | `false`                                                      | Enables ingress controller                                                       |
| ingress.className                             | string  | `""`                                                         | Ingress class name                                                               |
| ingress.annotations                           | object  | `{}`                                                         | Ingress annotations                                                              |
| ingress.hosts                                 | array   | `[]`                                                         | List of host rules                                                               |
| ingress.tls                                   | array   | `[]`                                                         | Ingress TLS configuration                                                        |
| resources                                     | object  | `{}`                                                         | Resource requirements for the relay container                                    |
| autoscaling.enabled                           | bool    | `false`                                                      | Enables HorizontalPodAutoscaler                                                  |
| autoscaling.minReplicas                       | integer | `1`                                                          | Sets minimum number of running replicas                                          |
| autoscaling.maxReplicas                       | integer | `100`                                                        | Sets maximum number of running replicas                                          |
| autoscaling.targetCPUUtilizationPercentage    | integer | `80`                                                         | Configures CPU as an average utilization metrics resource                        |
| autoscaling.targetMemoryUtilizationPercentage | integer | `80`                                                         | Configures memory as an average utilization metrics resource                     |
| nodeSelector                                  | object  | `{}`                                                         | Selector to target node placement for the relay pod                              |
| tolerations                                   | array   | `[]`                                                         | Specify pod tolerations                                                          |
| affinity                                      | object  | `{}`                                                         | Specify pod affinity                                                             |

## Learn more

Read our [documentation](https://docs.launchdarkly.com) for in-depth instructions on configuring and using LaunchDarkly. To learn more about the Relay Proxy specifically, read the [complete reference guide for the Relay Proxy](https://docs.launchdarkly.com/home/relay-proxy).

## Contributing

We encourage pull requests and other contributions from the community. Check out our [contributing guidelines](CONTRIBUTING.md) for instructions on how to contribute to this repository.

## About LaunchDarkly

* LaunchDarkly is a continuous delivery platform that provides feature flags as a service and allows developers to iterate quickly and safely. We allow you to easily flag your features and manage them from the LaunchDarkly dashboard.  With LaunchDarkly, you can:
    * Roll out a new feature to a subset of your users (like a group of users who opt-in to a beta tester group), gathering feedback and bug reports from real-world use cases.
    * Gradually roll out a feature to an increasing percentage of users, and track the effect that the feature has on key metrics (for instance, how likely is a user to complete a purchase if they have feature A versus feature B?).
    * Turn off a feature that you realize is causing performance problems in production, without needing to re-deploy, or even restart the application with a changed configuration file.
    * Grant access to certain features based on user attributes, like payment plan (eg: users on the ‘gold’ plan get access to more features than users in the ‘silver’ plan). Disable parts of your application to facilitate maintenance, without taking everything offline.
* LaunchDarkly provides feature flag SDKs for a wide variety of languages and technologies. Check out [our documentation](https://docs.launchdarkly.com/docs) for a complete list.
* Explore LaunchDarkly
    * [launchdarkly.com](https://www.launchdarkly.com/ "LaunchDarkly Main Website") for more information
    * [docs.launchdarkly.com](https://docs.launchdarkly.com/  "LaunchDarkly Documentation") for our documentation and SDK reference guides
    * [apidocs.launchdarkly.com](https://apidocs.launchdarkly.com/  "LaunchDarkly API Documentation") for our API documentation
    * [launchdarkly.com/blog](https://launchdarkly.com/blog/  "LaunchDarkly Blog Documentation") for the latest product updates
