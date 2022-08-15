# LaunchDarkly Relay Proxy Helm Chart

A helm chart to ease deployment of the [LaunchDarkly Relay Proxy](https://github.com/launchdarkly/ld-relay) to Kubernetes.

## LaunchDarkly overview

[LaunchDarkly](https://www.launchdarkly.com) is a feature management platform that serves over 100 billion feature flags daily to help teams build better software, faster. [Get started](https://docs.launchdarkly.com/docs/getting-started) using LaunchDarkly today!

[![Twitter Follow](https://img.shields.io/twitter/follow/launchdarkly.svg?style=social&label=Follow&maxAge=2592000)](https://twitter.com/intent/follow?screen_name=launchdarkly)

## Installation

```shell
helm repo add launchdarkly-ld-relay https://launchdarkly.github.io/ld-relay-helm
helm install ld-relay launchdarkly-ld-relay/ld-relay
```

This command will deploy the relay proxy to the Kubernetes cluster using the default configuration. The default configuration is insufficient to have a working instance of the proxy running. You must minimally provide an environment for the proxy to connect to. See the configuration section below.

## Configuration options

This chart can be customized by overriding the configuration options defined in [the values file](https://github.com/launchdarkly/ld-relay-helm/blob/main/values.yaml). You are encouraged to review this file as it contains detailed documentation. The list of value options are also summarized below.

The relay proxy is controlled through environment variables. These can be set directly by specifying a name and value in the `relay.environment` option, or through a secret using the `relay.secrets` option. See the [relay proxy's guide on configuration](https://github.com/launchdarkly/ld-relay/blob/v6/docs/configuration.md#file-section-environment-name) for a list of valid environment variable names

| Key                                             | Type    | Default                                                      | Description                                                                    |
| ----------------------------------------------- | ------- | ------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| relay.environment                               | object  | `{}`                                                         | Define container environment variables to configure the relay instance         |
| relay.secrets                                   | array   | `[]`                                                         | Define container environment variables populated from a k8s secret             |
| replicaCount                                    | integer | `1`                                                          | Number of replicas of the relay pod                                            |
| image.repository                                | string  | `launchdarkly/ld-relay`                                      | ld-relay image repository                                                      |
| image.pullPolicy                                | string  | `IfNotPresent`                                               | ld-relay image pull policy                                                     |
| image.tag                                       | string  | `""`                                                         | Overrides the image tag whose default is the chart appVersion.                 |
| imagePullSecrets                                | array   | `[]`                                                         | Specify dockere registry secret names as an array                              |
| nameOverride                                    | string  | `""`                                                         | Partially override the fullname template with a string (includes release name) |
| fullnameOverride                                | string  | `""`                                                         | Fully override the fullname template with a string                             |
| serviceAccount.create                           | bool    | `true`                                                       | Specifies whether a service account should be created                          |
| serviceAccount.annotations                      | object  | `{}`                                                         | Annotations to add to the service account                                      |
| serviceAccount.name                             | string  | `""`                                                         | The name of the service account to use.                                        |
| podAnnotations                                  | object  | `{}`                                                         | Pod annotations                                                                |
| podSecurityContext                              | object  | `{}`                                                         | Pod security context                                                           |
| securityContext                                 | object  | `{}`                                                         | Container security context                                                     |
| service.type                                    | string  | `ClusterIP`                                                  | Kubernetes service type                                                        |
| service.ports                                   | array   | `[{port: 8030, targetPort: 8030, protocol: TCP, name: api}]` | Service port mapping. Must include one port named api.                         |
| ingress.enabled                                 | bool    | `false`                                                      | Enable ingress controller                                                      |
| ingress.className                               | string  | `""`                                                         | Ingress class name                                                             |
| ingress.annotations                             | object  | `{}`                                                         | Ingress annotations                                                            |
| ingress.hosts                                   | array   | `[]`                                                         | List of host rules                                                             |
| ingress.tls                                     | array   | `[]`                                                         | Ingress TLS configuration                                                      |
| resources                                       | object  | `{}`                                                         | Resource requirements for the relay container                                  |
| autoscaling.enabled                             | bool    | `false`                                                      | Enable HorizontalPodAutoscaler                                                 |
| autoscaling.minReplicas                         | integer | `1`                                                          | Set minimum number of running replicas                                         |
| autoscaling.maxReplicas                         | integer | `100`                                                        | Set maximum number of running replicas                                         |
| autoscaling.targetCPUUtilizationPercentage      | integer | `80`                                                         | Configure CPU as an average utilization metrics resource                       |
| autoscaling.targetMemoryUtilizationPercentage   | integer | `80`                                                         | Configure memory as an average utilization metrics resource                    |
| nodeSelector                                    | object  | `{}`                                                         | Selector to target node placement for the relay pod                            |
| tolerations                                     | array   | `[]`                                                         | Specify pod tolerations                                                        |
| affinity                                        | object  | `{}`                                                         | Specify pod affinity                                                           |

### Examples

There are multiple ways to override the chart values shown above. For the purposes of these examples, we will show an override.yaml file which can be used like:

```shell
helm install --values ./override.yaml ld-relay launchdarkly-ld-relay/ld-relay
```

**Minimal example**

```yaml
# override.yaml
relay:
  environment:
    LD_ENV_Production: your-sdk-key
```

**Configuring using secrets**


```yaml
relay:
  # Specify the relay environment variables here to load them into this chart's ConfigMap directly.
  environment:
    USE_REDIS: true
    REDIS_HOST: redis-master.default.svc.cluster.local
  secrets:
    - envName: LD_ENV_Production
      secretName: relay
      secretKey: sdk-key
    - envName: REDIS_PASSWORD
      secretName: relay
      secretKey: redis-password
```

In the above example, both the SDK key and the password for the redis cluster are being pulled from pre-existing Kubernetes secrets. e.g.,

```shell
kubectl create secret generic relay --from-literal=redis-password=your-password --from-literal=sdk-key=your-sdk-key
```

## Learn more

Check out our [documentation](https://docs.launchdarkly.com) for in-depth instructions on configuring and using LaunchDarkly. You can also head straight to the [complete reference guide for the relay proxy](https://docs.launchdarkly.com/home/relay-proxy).

## Testing

The unit tests start with a base set of comparison files. These are referred to as "golden files". These should
represent an identical match to the installed chart if only the default values are used. Starting from this base point
allows us to force additional unit tests only on the properties we expect to change.

Golden files can be updated by running `make update-golden-files`.
Unit tests are run with `make test`.

> **NOTE** The testing structure drew heavy inspiration from [this blog post](https://camunda.com/blog/2022/03/test/).

## Contributing

We encourage pull requests and other contributions from the community.

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
