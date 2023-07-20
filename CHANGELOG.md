Change log
================================================

All notable changes to the LaunchDarkly Relay Proxy Helm Chart will be documented in this file. This project adheres to [Semantic Versioning](https://semver.org).


## [2.4.0] - 2023-07-20
### Added:
- Support TopologySpreadConstraint and PodDistruptionBudget. (Thanks, [pjaak](https://github.com/launchdarkly/ld-relay-helm/pull/47)!)

## [2.3.0] - 2023-07-04
### Added:
- Support Dynamic EnvFrom existing K8s Secrets. (Thanks, [LiamStorkey](https://github.com/launchdarkly/ld-relay-helm/pull/45)!)

## [2.2.2] - 2023-05-11
### Fixed:
- Liveness and readiness probes, by default use an HTTP scheme. This breaks once TLS is enabled in the relay. values file now updated to allow these schemes to be overridden.

## [2.2.1] - 2023-05-10
### Fixed:
- Fix name collision when setting multiple secret values as volume mounts.

### Removed:
- `relay.secrets.volumeName` is no longer used when mounting secrets as volumes as all secrets as mounted within a shared volume.

## [2.2.0] - 2023-04-27
### Added:
- Support setting environment variables directly on the container spec through `relay.environmentVariables`. This enables more complex variable definitions. (Thanks, [uristernik](https://github.com/launchdarkly/ld-relay-helm/pull/34)!)

## [2.1.0] - 2023-04-05
### Added:
- Add ability to set labels on pods through `pod.labels` value. (Thanks, [kovaxur](https://github.com/launchdarkly/ld-relay-helm/pull/30)!)

### Deprecated:
- `podAnnotations` and `podSecurityContext` values have been deprecated. Use `pod.annotations` and `pod.securityContext` instead.

## [2.0.0] - 2023-03-24
### Changed:
- Updated default Relay Proxy version to v7.2.1 to support [contexts](https://docs.launchdarkly.com/home/contexts).

## [1.2.1] - 2023-02-24
### Changed:
- (Tests) Bumped golang.org/x/text from 0.3.6 to 0.3.8
- (Tests) Bumped golang.org/x/net from 0.0.0-20210614182718-04defd469f4e to 0.7.0
- (Tests) Bumped golang.org/x/crypto from 0.0.0-20210513164829-c07d793c2f9a to 0.1.0

## [1.2.0] - 2023-01-25
### Changed:
- Updated HorizontalPodAutoscaler to be compatible with older versions of Kubernetes. (Thanks, [guifl](https://github.com/launchdarkly/ld-relay-helm/pull/21)!)

## [1.1.0] - 2022-11-16
### Added:
- Allow setting annotations on the created service.
- Add mechanism for mounting secrets as files.

## [1.0.1] - 2022-11-07
### Fixed:
- Restart running containers if the ConfigMap values change.

## [1.0.0] - 2022-10-14
### Added:
- Initial release of the LaunchDarkly Relay Proxy Helm Chart
