# LaunchDarkly Relay Proxy Helm Chart

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/ld-relay)](https://artifacthub.io/packages/search?repo=ld-relay)
[![Quality control](https://github.com/launchdarkly/ld-relay-helm/actions/workflows/ci.yml/badge.svg)](https://github.com/launchdarkly/ld-relay-helm/actions/workflows/ci.yml)

A Helm chart to ease deployment of the [LaunchDarkly Relay Proxy](https://github.com/launchdarkly/ld-relay) to Kubernetes (k8s).

Basic [installation](#installation) and [configuration](#configuration-options) information is below.

To learn more, read the [Getting started](./docs/getting-started.md) guide. For additional examples, read:

- [Offline mode](./docs/examples/offline-mode.md)
- [Automatic configuration](./docs/examples/automatic-configuration.md)

## LaunchDarkly overview

[LaunchDarkly](https://www.launchdarkly.com) is a feature management platform that serves trillions of feature flags daily to help teams build better software, faster. [Get started](https://docs.launchdarkly.com/docs/getting-started) using LaunchDarkly today!

[![Twitter Follow](https://img.shields.io/twitter/follow/launchdarkly.svg?style=social&label=Follow&maxAge=2592000)](https://twitter.com/intent/follow?screen_name=launchdarkly)

## Installation

The default configuration is insufficient to have a working instance of the Relay Proxy running. You must minimally provide an environment for the Relay Proxy to connect to, using your LaunchDarkly SDK key for that environment.

To deploy the Relay Proxy to the Kubernetes cluster using the default configuration and an environment:

```shell
helm repo add launchdarkly-ld-relay https://launchdarkly.github.io/ld-relay-helm
helm install ld-relay --set relay.environment.LD_ENV_YourEnvironment=your-sdk-key launchdarkly-ld-relay/ld-relay
```

For additional configuration, use the [Configuration options](#configuration-options) below.

OCI artifacts of all LaunchDarkly Helm charts are available in [ghcr.io](https://github.com/orgs/launchdarkly/packages?q=helm-charts).

## Kubernetes version support

This chart follows the Kubernetes project's [support policy](https://kubernetes.io/releases/version-skew-policy/#supported-versions). We test the chart against the three newest Kubernetes minor versions and fix bugs that affect them. Each minor version gets this support for about 14 months, as it does upstream. The chart can work on older versions, but we do not test them or fix bugs that occur only on them.

The chart's `kubeVersion` setting is a separate, lower floor. Helm refuses to install or upgrade the chart on a cluster older than this floor. The floor is the oldest Kubernetes version that Amazon EKS or Google GKE still supports, including their paid extended support. This keeps the chart installable on clusters that your provider still supports. We raise the floor only in a major chart release.

| Chart version | Minimum Kubernetes version |
|---------------|----------------------------|
| 4.x           | 1.32                       |

Without a cluster, `helm template` checks the floor against a default Kubernetes version built into Helm. If that default is older than the floor, pass your cluster's version with `--kube-version`.

## Security defaults

The chart's default `pod.securityContext` and `securityContext` satisfy the Kubernetes [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/) `restricted` profile. The pod uses the `RuntimeDefault` seccomp profile. The Relay Proxy container sets `runAsNonRoot`, sets `allowPrivilegeEscalation: false` and drops all capabilities.

The chart does not set `runAsUser`, `runAsGroup` or `fsGroup`, so the image decides which user the Relay Proxy runs as. This also lets platforms that assign UIDs, such as OpenShift, admit the pod.

Helm merges your values with these defaults:

- Setting `securityContext: {}` does not remove them.
- To remove one default, set it to `null`, for example `securityContext.runAsNonRoot: null`.
- To remove all of them, set both `securityContext: null` and `pod.securityContext: null`. This gives the behaviour of chart versions before 4.0.
- Lists replace the default rather than merge with it. For example, `securityContext.capabilities.drop: [NET_RAW]` no longer drops `ALL`.
- When this chart is a subchart, `null` in the parent chart's `values.yaml` does not remove a default. Pass the `null` with `-f` or `--set` instead.

These container-level defaults do not apply to `relay.initContainers`. In a namespace that enforces `restricted`, give each init container its own `securityContext` with the same fields.

To upgrade from a chart version before 4.0, see [Upgrading to 4.0](#upgrading-to-40).

To harden the pod further:

- Set `securityContext.readOnlyRootFilesystem: true`. If you also use [offline mode](docs/examples/offline-mode.md), mount a writable volume at `/tmp`, because offline mode expands its archive there. Without it, the Relay Proxy does not start.

    ```yaml
    securityContext:
      readOnlyRootFilesystem: true
    relay:
      extraVolumes:
        - name: tmp
          emptyDir: {}
      extraVolumeMounts:
        - name: tmp
          mountPath: /tmp
    ```

- Set `automountServiceAccountToken: false`. The Relay Proxy does not call the Kubernetes API. Keep the token if another part of your setup reads it, for example Vault Agent with Kubernetes auth.

## Configuration options

To customize this Helm chart, override the configuration options defined in the [values file](https://github.com/launchdarkly/ld-relay-helm/blob/main/values.yaml). The values file contains detailed documentation on each option.

Here's a summary of the available configuration options:

| Key                                           | Type    | Default                                                      | Description                                                                                              |
| --------------------------------------------- | ------- | ------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------- |
| relay.environment                             | object  | `{}`                                                         | Defines container environment variables to configure the Relay Proxy instance (via ConfigMap)            |
| relay.environmentVariables                    | array   | `[]`                                                         | Defines container environment variables to configure the Relay Proxy instance (via container spec)       |
| relay.envFromSecret                           | string  | `null`                                                       | Defines container environment variables to configure the Relay Proxy instance (via existing k8s secrets) |
| relay.secrets                                 | array   | `[]`                                                         | Defines container environment variables or volumes populated from k8s secrets                            |
| relay.volume                                  | object  | `{}`                                                         | Enables offline mode or references an existing config file from a defined volume                         |
| relay.livenessProbe                           | object  | `{httpGet: { port: "api",  path: "/status" }}`               | Defines the liveness probe for the relay container                                                       |
| relay.readinessProbe                          | object  | `{httpGet: { port: "api",  path: "/status?expect=..." }}`    | Defines the readiness probe for the relay container. The default `expect` clauses pass once every environment has initialized |
| relay.startupProbe                            | object  | `null`                                                       | Defines the startup probe for the relay  container                                                     |
| relay.lifecycle                               | object  | `{}`                                                         | Defines the lifecycle hooks for the relay container                                                      |
| replicaCount                                  | integer | `1`                                                          | Number of replicas of the relay pod                                                                      |
| image.repository                              | string  | `launchdarkly/ld-relay`                                      | ld-relay image repository                                                                                |
| image.pullPolicy                              | string  | `IfNotPresent`                                               | ld-relay image pull policy                                                                               |
| image.tag                                     | string  | `""`                                                         | Overrides the image tag whose default is the chart appVersion                                            |
| imagePullSecrets                              | array   | `[]`                                                         | Specifies docker registry secret names as an array                                                       |
| nameOverride                                  | string  | `""`                                                         | Partially overrides the fullname template with a string (includes release name)                          |
| fullnameOverride                              | string  | `""`                                                         | Fully overrides the fullname template with a string                                                      |
| namespaceOverride                             | string  | `""`                                                         | Sets `metadata.namespace` on namespaced resources; defaults to the release namespace. Refer to `values.yaml` for when to use it |
| commonLabels                                  | object  | `{}`                                                         | Specifies labels to be applied to all created resources                                                  |
| serviceAccount.create                         | bool    | `true`                                                       | Specifies whether a service account should be created                                                    |
| serviceAccount.annotations                    | object  | `{}`                                                         | Annotations to add to the service account                                                                |
| serviceAccount.name                           | string  | `""`                                                         | The name of the service account                                                                          |
| automountServiceAccountToken                  | bool    | `null`                                                       | Pod automountServiceAccountToken. When unset, Kubernetes applies its default                             |
| pod.annotations                               | object  | `{}`                                                         | Pod annotations                                                                                          |
| pod.labels                                    | object  | `{}`                                                         | Pod labels                                                                                               |
| pod.securityContext                           | object  | `{seccompProfile: {type: RuntimeDefault}}`                   | Pod security context                                                                                     |
| terminationGracePeriodSeconds                 | string  | `null`                                                       | Pod terminationGracePeriodSeconds                                                                        |
| securityContext                               | object  | `{allowPrivilegeEscalation: false, capabilities: {drop: [ALL]}, runAsNonRoot: true, seccompProfile: {type: RuntimeDefault}}` | Container security context                                                                               |
| service.type                                  | string  | `ClusterIP`                                                  | Kubernetes service type                                                                                  |
| service.annotations                           | object  | `{}`                                                         | Annotations to add to the service                                                                        |
| service.ports                                 | array   | `[{port: 8030, targetPort: 8030, protocol: TCP, name: api}]` | Service port mapping. Must include one port named `api`.                                                 |
| ingress.enabled                               | bool    | `false`                                                      | Enables ingress controller                                                                               |
| ingress.className                             | string  | `""`                                                         | Ingress class name                                                                                       |
| ingress.annotations                           | object  | `{}`                                                         | Ingress annotations                                                                                      |
| ingress.hosts                                 | array   | `[]`                                                         | List of host rules                                                                                       |
| ingress.tls                                   | array   | `[]`                                                         | Ingress TLS configuration                                                                                |
| httpRoute.enabled                             | bool    | `false`                                                      | Enables Gateway API HTTPRoute                                                                            |
| httpRoute.annotations                         | object  | `{}`                                                         | HTTPRoute annotations                                                                                    |
| httpRoute.parentRefs                          | array   | `[]`                                                         | List of Gateway parent references this route attaches to (full ParentReference schema; rendered as-is)   |
| httpRoute.hostnames                           | array   | `[]`                                                         | List of hostnames this route matches                                                                     |
| httpRoute.rules                               | array   | `[]`                                                         | List of routing rules; see the fields below                                                              |
| httpRoute.rules[].matches                     | array   | `[]`                                                         | Full Gateway API match schema (path, headers, method, queryParams); rendered as-is                       |
| httpRoute.rules[].filters                     | array   | `[]`                                                         | Full Gateway API filter schema (header modifiers, redirect, URL rewrite, mirror); rendered as-is         |
| httpRoute.rules[].timeouts                    | object  | `{}`                                                         | Gateway API timeouts (request, backendRequest); rendered as-is                                           |
| httpRoute.rules[].port                        | integer | `nil`                                                        | Convenience: routes the rule's backend to the chart's Service on this port                               |
| httpRoute.rules[].backendRefs                 | array   | `[]`                                                         | Explicit backends (weighting, mirroring, external); overrides the `port` convenience when set            |
| resources                                     | object  | `{}`                                                         | Resource requirements for the relay container                                                            |
| autoscaling.enabled                           | bool    | `false`                                                      | Enables HorizontalPodAutoscaler                                                                          |
| autoscaling.minReplicas                       | integer | `1`                                                          | Sets minimum number of running replicas                                                                  |
| autoscaling.maxReplicas                       | integer | `100`                                                        | Sets maximum number of running replicas                                                                  |
| autoscaling.targetCPUUtilizationPercentage    | integer | `80`                                                         | Configures CPU as an average utilization metrics resource                                                |
| autoscaling.targetMemoryUtilizationPercentage | integer | `80`                                                         | Configures memory as an average utilization metrics resource                                             |
| autoscaling.behavior                          | object  | `null`                                                       | Optional HPA scale behavior (requires Kubernetes >= 1.18)                                                |
| autoscaling.additionalMetrics                 | array   | `[]`                                                         | Extra External/Pods HPA metrics appended verbatim to `spec.metrics` (requires Kubernetes >= 1.12)        |
| deployment.strategy                           | object  | `null`                                                       | Optional rolling update strategy for the deployment                                                      |
| deployment.revisionHistoryLimit               | integer | `null`                                                       | Optional number of old ReplicaSets to retain for rollback (Kubernetes default `10`)                      |
| nodeSelector                                  | object  | `{}`                                                         | Selector to target node placement for the relay pod                                                      |
| tolerations                                   | array   | `[]`                                                         | Specify pod tolerations                                                                                  |
| affinity                                      | object  | `{}`                                                         | Specify pod affinity                                                                                     |
| pod.disruptionBudget.enabled                  | boolean | `false`                                                      | Enabled podDisruptionBudget                                                                              |
| pod.disruptionBudget.minAvailable             | string  | `""`                                                         | Minimum number of pods that are available after eviction as number or percentage                         |
| pod.disruptionBudget.maxUnavailable           | string  | `""`                                                         | Maximum number of pods that are unavailable after eviction as number or percentage                       |
| pod.topologySpreadConstraints                 | array   | `[]`                                                         | Specify the topology spread constrait definitions to apply to the relay deployment                       |
| pod.priorityClassName                         | string  | `""`                                                         | Specify a PriorityClass for the pod                                                                      |
| pod.dnsPolicy                                 | string  | `""`                                                         | Optional pod DNS policy (Kubernetes default `ClusterFirst`); set to `None` when supplying a custom dnsConfig |
| pod.dnsConfig                                 | object  | `{}`                                                         | Optional pod DNS configuration (nameservers, searches, options)                                          |

`httpRoute` requires the [Gateway API](https://gateway-api.sigs.k8s.io/) CRDs (`gateway.networking.k8s.io`) to already be installed in the cluster, along with a `Gateway` resource for the route to attach to. Use either `ingress` or `httpRoute`, not both.

For each rule, set `port` to route to this chart's Service, or set `backendRefs` explicitly. `backendRefs` takes precedence over `port` when both are given, and (unlike `port`) can target any Service, including another namespace. A rule with no `matches` matches every request and publishes the entire relay -- including admin/metrics endpoints -- on the route's hostname, so add explicit `matches` unless a catch-all is intended.

## Upgrading to 4.0

### Relay Proxy v9

Chart 4.0 installs Relay Proxy v9 by default. Read the [Relay Proxy v9 changelog](https://github.com/launchdarkly/ld-relay/blob/v9/CHANGELOG.md) for the changes to Relay itself, such as the move from OpenCensus to OpenTelemetry. The chart changes are:

- The relay image is distroless. It has no shell, so `exec` probes and lifecycle hooks that run `sh` no longer work. Use `httpGet` probes, and use the `sleep` action for a `preStop` delay (Kubernetes 1.30 or later). The `-static-debian13-debug-nonroot` image tags include a busybox shell at `/busybox/sh` for debugging.
- The image's user changed from uid 1000 to the distroless nonroot user (uid 65532). If you set `runAsUser`, `runAsGroup` or `fsGroup` to 1000, remove those overrides, or make sure the files in your volumes are readable by the user you choose.
- The default readiness probe uses `/status` `expect` assertions. A pod becomes ready when it has at least one environment and every environment has finished initializing. It stays ready during a LaunchDarkly outage. The liveness probe still requests the plain `/status` endpoint.
- The Prometheus exporter was removed from Relay v9, so the values file no longer shows a Prometheus port. Export metrics over OTLP instead.
- The deprecated `pod.distruptionBudget` value is removed. Use `pod.disruptionBudget`.

### Kubernetes version

Chart 4.0 requires Kubernetes 1.32 or later. The templates no longer produce the older API versions of the HorizontalPodAutoscaler, PodDisruptionBudget and Ingress resources. On older clusters, stay on chart 3.x. For the policy, see [Kubernetes version support](#kubernetes-version-support).

### Security context

Chart versions before 4.0 set no security context. When you upgrade, check these cases:

- **Custom images.** The image must run as a numeric, non-root UID. An image that runs as root, or that sets `USER` to a name, fails the `runAsNonRoot` check and the container does not start. The published `launchdarkly/ld-relay` images run as a numeric non-root UID.
- **Settings that conflict with the defaults.** Your values merge with the defaults, so some settings now conflict. `securityContext.runAsUser: 0` conflicts with `runAsNonRoot: true`. `securityContext.privileged: true` and capabilities in `securityContext.capabilities.add` such as `SYS_ADMIN` conflict with `allowPrivilegeEscalation: false`. Remove the conflicting default with `null`.
- **Pod-level seccomp profiles.** The container now sets `seccompProfile: RuntimeDefault`, which takes precedence over `pod.securityContext.seccompProfile`. If you set a different profile on the pod, set it on `securityContext.seccompProfile` as well.
- **`helm upgrade --reuse-values`.** This keeps the values of your previous release, including its empty security contexts, so the new defaults do not apply. Use `--reset-then-reuse-values` (Helm 3.14 or later) or pass your values explicitly.

## Learn more

Read our [documentation](https://docs.launchdarkly.com) for in-depth instructions on configuring and using LaunchDarkly. To learn more about the Relay Proxy specifically, read the [complete reference guide for the Relay Proxy](https://docs.launchdarkly.com/home/relay-proxy).

## Contributing

We encourage pull requests and other contributions from the community. Check out our [contributing guidelines](CONTRIBUTING.md) for instructions on how to contribute to this repository.

## About LaunchDarkly

- LaunchDarkly is a continuous delivery platform that provides feature flags as a service and allows developers to iterate quickly and safely. We allow you to easily flag your features and manage them from the LaunchDarkly dashboard. With LaunchDarkly, you can:
  - Roll out a new feature to a subset of your users (like a group of users who opt-in to a beta tester group), gathering feedback and bug reports from real-world use cases.
  - Gradually roll out a feature to an increasing percentage of users, and track the effect that the feature has on key metrics (for instance, how likely is a user to complete a purchase if they have feature A versus feature B?).
  - Turn off a feature that you realize is causing performance problems in production, without needing to re-deploy, or even restart the application with a changed configuration file.
  - Grant access to certain features based on user attributes, like payment plan (eg: users on the ‘gold’ plan get access to more features than users in the ‘silver’ plan). Disable parts of your application to facilitate maintenance, without taking everything offline.
- LaunchDarkly provides feature flag SDKs for a wide variety of languages and technologies. Read [our documentation](https://docs.launchdarkly.com/docs) for a complete list.
- Explore LaunchDarkly
  - [launchdarkly.com](https://www.launchdarkly.com/ "LaunchDarkly Main Website") for more information
  - [docs.launchdarkly.com](https://docs.launchdarkly.com/ "LaunchDarkly Documentation") for our documentation and SDK reference guides
  - [apidocs.launchdarkly.com](https://apidocs.launchdarkly.com/ "LaunchDarkly API Documentation") for our API documentation
  - [launchdarkly.com/blog](https://launchdarkly.com/blog/ "LaunchDarkly Blog Documentation") for the latest product updates
