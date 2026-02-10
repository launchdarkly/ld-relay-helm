Change log
================================================

All notable changes to the LaunchDarkly Relay Proxy Helm Chart will be documented in this file. This project adheres to [Semantic Versioning](https://semver.org).


## [3.5.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.4.1...3.5.0) (2026-02-10)


### Features

* Add support for init containers and multiple volumes ([#90](https://github.com/launchdarkly/ld-relay-helm/issues/90)) ([c5cd4ff](https://github.com/launchdarkly/ld-relay-helm/commit/c5cd4ff999d3210d2a613ca7e21dc166773f4cce))

## [3.4.1](https://github.com/launchdarkly/ld-relay-helm/compare/3.4.0...3.4.1) (2026-01-30)


### Bug Fixes

* Use fullname for config resources ([#86](https://github.com/launchdarkly/ld-relay-helm/issues/86)) ([b71f668](https://github.com/launchdarkly/ld-relay-helm/commit/b71f6687e6868fde9dceeaff87d42b7a540bf334)), closes [#85](https://github.com/launchdarkly/ld-relay-helm/issues/85)

## [3.4.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.3.2...3.4.0) (2025-05-20)


### Features

* Add commonLabels property ([#77](https://github.com/launchdarkly/ld-relay-helm/issues/77)) ([fd3acfe](https://github.com/launchdarkly/ld-relay-helm/commit/fd3acfe7af1155488ca669f4c28713d7d4dc1302))

## [3.3.2](https://github.com/launchdarkly/ld-relay-helm/compare/3.3.1...3.3.2) (2025-05-05)


### Bug Fixes

* Correct typo for PodDisruptionBudget ([#75](https://github.com/launchdarkly/ld-relay-helm/issues/75)) ([3042465](https://github.com/launchdarkly/ld-relay-helm/commit/304246580a75e4cda0b5fdbbfb922b1d621a1204))

## [3.3.1](https://github.com/launchdarkly/ld-relay-helm/compare/3.3.0...3.3.1) (2025-01-16)


### Bug Fixes

* Bump default relay proxy version from 8.2.0 to 8.10.5 ([#69](https://github.com/launchdarkly/ld-relay-helm/issues/69)) ([fc4693b](https://github.com/launchdarkly/ld-relay-helm/commit/fc4693babd4db37ba36ffe41ebcea1a72a1f68f6))

## [3.3.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.2.0...3.3.0) (2024-03-19)


### Features

* Add option to define terminationGracePeriodSeconds ([#60](https://github.com/launchdarkly/ld-relay-helm/issues/60)) ([f65b60b](https://github.com/launchdarkly/ld-relay-helm/commit/f65b60b0cc0f1f956e3a951042095583e85cd542))

## [3.2.0] - 2024-02-21
### Added:
- Add support for container lifecycle hooks. (Thanks, [Helinanu](https://github.com/launchdarkly/ld-relay-helm/pull/57)!)

## [3.1.0] - 2023-12-01
### Added:
- Support setting a pod's priority class name. (Thanks, [kh3dron](https://github.com/launchdarkly/ld-relay-helm/pull/53)!)

## [3.0.0] - 2023-10-25
### Changed:
- Updated the default relay image to v8.2.0. To learn more about the changes involved, read [the Relay Changelog](https://github.com/launchdarkly/ld-relay/blob/v8/CHANGELOG.md).

### Removed:
- Removed previously deprecated config option `podAnnotations`. Use `pod.annotations` instead.
- Removed previously deprecated config option `podSecurityContext`. Use `pod.securityContext` instead.

## [2.4.0] - 2023-07-20
### Added:
- Support TopologySpreadConstraint and PodDisruptionBudget. (Thanks, [pjaak](https://github.com/launchdarkly/ld-relay-helm/pull/47)!)

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
