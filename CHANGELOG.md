Change log
================================================

All notable changes to the LaunchDarkly Relay Proxy Helm Chart will be documented in this file. This project adheres to [Semantic Versioning](https://semver.org).


## [3.12.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.11.4...3.12.0) (2026-11-19)


### Features

* Add commonLabels property ([#77](https://github.com/launchdarkly/ld-relay-helm/issues/77)) ([fd3acfe](https://github.com/launchdarkly/ld-relay-helm/commit/fd3acfe7af1155488ca669f4c28713d7d4dc1302))
* Add explicit namespace support ([#115](https://github.com/launchdarkly/ld-relay-helm/issues/115)) ([90f85fd](https://github.com/launchdarkly/ld-relay-helm/commit/90f85fdb2d4b16cb9d502a2ff19e1a5719a005a2))
* Add option to configure startup probes in ld-relay container ([#103](https://github.com/launchdarkly/ld-relay-helm/issues/103)) ([cdac3d5](https://github.com/launchdarkly/ld-relay-helm/commit/cdac3d5eb39657e46282c5ab1b7511ac590a7649))
* Add option to define terminationGracePeriodSeconds ([#60](https://github.com/launchdarkly/ld-relay-helm/issues/60)) ([f65b60b](https://github.com/launchdarkly/ld-relay-helm/commit/f65b60b0cc0f1f956e3a951042095583e85cd542))
* Add pod priorityClassName support ([#53](https://github.com/launchdarkly/ld-relay-helm/issues/53)) ([e28b01b](https://github.com/launchdarkly/ld-relay-helm/commit/e28b01b1b7437332bbb73a522f7b739cb670809d))
* Add support for container lifecycle hooks ([#57](https://github.com/launchdarkly/ld-relay-helm/issues/57)) ([f47acc6](https://github.com/launchdarkly/ld-relay-helm/commit/f47acc6a94a76c0550dfce2dfb4d5d793bdc9050))
* Add support for HTTP routing ([#132](https://github.com/launchdarkly/ld-relay-helm/issues/132)) ([5c79d98](https://github.com/launchdarkly/ld-relay-helm/commit/5c79d9885941a4732482f536b8ebd8a6770b81c7))
* Add support for init containers and multiple volumes ([#90](https://github.com/launchdarkly/ld-relay-helm/issues/90)) ([c5cd4ff](https://github.com/launchdarkly/ld-relay-helm/commit/c5cd4ff999d3210d2a613ca7e21dc166773f4cce))
* Add support for user-defined pod labels ([#30](https://github.com/launchdarkly/ld-relay-helm/issues/30)) ([99e105d](https://github.com/launchdarkly/ld-relay-helm/commit/99e105d998efdb3dde9d6a3448e2b642c111cce8))
* Add TopologySpreadConstraint and PodDistruptionBudget support ([#47](https://github.com/launchdarkly/ld-relay-helm/issues/47)) ([9bdecf9](https://github.com/launchdarkly/ld-relay-helm/commit/9bdecf9d0a4ef08c94b81850c0cfd81af7b4146a))
* allow additional HPA metrics via autoscaling.additionalMetrics ([#143](https://github.com/launchdarkly/ld-relay-helm/issues/143)) ([ba05727](https://github.com/launchdarkly/ld-relay-helm/commit/ba0572711d584ee38bf1c6bdfb7efa5cfc033d8d))
* Default to 7.2.1 version of the relay proxy ([#28](https://github.com/launchdarkly/ld-relay-helm/issues/28)) ([7cd59d5](https://github.com/launchdarkly/ld-relay-helm/commit/7cd59d505e93585fd8e8bcb35ae15d7afb2fec60))
* expose deployment rolling update strategy and HPA scale behavior as configurable values ([#117](https://github.com/launchdarkly/ld-relay-helm/issues/117)) ([c3b2259](https://github.com/launchdarkly/ld-relay-helm/commit/c3b22590d9f6c9549c04c9d4128c8b214a30c5b1))
* expose revisionHistoryLimit and dnsConfig/dnsPolicy as configurable values ([#125](https://github.com/launchdarkly/ld-relay-helm/issues/125)) ([c92136f](https://github.com/launchdarkly/ld-relay-helm/commit/c92136f2c59668295cfc6c2e1c2706e5bbe1cf58))
* Introduce top level pod configuration value ([#31](https://github.com/launchdarkly/ld-relay-helm/issues/31)) ([beb0f9e](https://github.com/launchdarkly/ld-relay-helm/commit/beb0f9efae74be389b40c618275d4138adcaab96))
* Remove deprecated pod values ([#52](https://github.com/launchdarkly/ld-relay-helm/issues/52)) ([444fc60](https://github.com/launchdarkly/ld-relay-helm/commit/444fc6079a053bd6e218c2ca8a85b3c58e89888b))
* Support Dynamic EnvFrom existing K8s Secrets ([#45](https://github.com/launchdarkly/ld-relay-helm/issues/45)) ([ae145e8](https://github.com/launchdarkly/ld-relay-helm/commit/ae145e8d7d1460baafc06a6e10c51e05c360a9f4))
* Support environment variables with reference values ([#35](https://github.com/launchdarkly/ld-relay-helm/issues/35)) ([7d9c1e6](https://github.com/launchdarkly/ld-relay-helm/commit/7d9c1e63c8a97e3c20891c9aa486d54c1eb54ad1))
* Update default relay to v8.2.0 ([#51](https://github.com/launchdarkly/ld-relay-helm/issues/51)) ([c82d92d](https://github.com/launchdarkly/ld-relay-helm/commit/c82d92d65754939d1609a696eba9adfc45aa17ee))


### Bug Fixes

* Allow multiple secrets to be mounted as volumes ([#37](https://github.com/launchdarkly/ld-relay-helm/issues/37)) ([4e9cd5c](https://github.com/launchdarkly/ld-relay-helm/commit/4e9cd5c19f30af59a995f79a5cdb0862d5a1b3ea))
* Allow overriding scheme for container probes ([#39](https://github.com/launchdarkly/ld-relay-helm/issues/39)) ([846d783](https://github.com/launchdarkly/ld-relay-helm/commit/846d7835103c4ed9a571f2ce050f6492484e4055))
* Bump default relay proxy version from 8.10.5 to 8.19.1 ([#112](https://github.com/launchdarkly/ld-relay-helm/issues/112)) ([b0eecc9](https://github.com/launchdarkly/ld-relay-helm/commit/b0eecc921d6d6ca93950db50ac99564ac152a24b))
* Bump default relay proxy version from 8.19.1 to 8.21.0 ([#139](https://github.com/launchdarkly/ld-relay-helm/issues/139)) ([348e242](https://github.com/launchdarkly/ld-relay-helm/commit/348e242ce3c20616c8d6e34c044d8629e6d2e4d9))
* Bump default relay proxy version from 8.2.0 to 8.10.5 ([#69](https://github.com/launchdarkly/ld-relay-helm/issues/69)) ([fc4693b](https://github.com/launchdarkly/ld-relay-helm/commit/fc4693babd4db37ba36ffe41ebcea1a72a1f68f6))
* Bump default relay proxy version from 8.21.0 to 8.22.0 ([#148](https://github.com/launchdarkly/ld-relay-helm/issues/148)) ([56ea574](https://github.com/launchdarkly/ld-relay-helm/commit/56ea574bb1faae4d44841219677b894d087dbb57))
* Bump default relay proxy version from 8.23.1 to 8.23.2 ([#154](https://github.com/launchdarkly/ld-relay-helm/issues/154)) ([3f9011c](https://github.com/launchdarkly/ld-relay-helm/commit/3f9011cff8e0948b50d3e46733ba79c3876802f1))
* Bump default relay proxy version from 8.23.2 to 8.23.3 ([#156](https://github.com/launchdarkly/ld-relay-helm/issues/156)) ([bda8a41](https://github.com/launchdarkly/ld-relay-helm/commit/bda8a4140192282748f8e7d40a7d7b6abd6285be))
* Bump default relay proxy version from 8.23.3 to 8.23.4 ([#158](https://github.com/launchdarkly/ld-relay-helm/issues/158)) ([194827a](https://github.com/launchdarkly/ld-relay-helm/commit/194827ac3a58814b1f4233d3fcc7f5ce7abb0bb9))
* Correct typo for PodDisruptionBudget ([#75](https://github.com/launchdarkly/ld-relay-helm/issues/75)) ([3042465](https://github.com/launchdarkly/ld-relay-helm/commit/304246580a75e4cda0b5fdbbfb922b1d621a1204))
* Default relay.lifecycle to map instead of array ([#99](https://github.com/launchdarkly/ld-relay-helm/issues/99)) ([613e261](https://github.com/launchdarkly/ld-relay-helm/commit/613e261bc287cde3951dc82a6e1afa8a0f9575ce))
* hash only the config data in the checksum/config annotation ([#146](https://github.com/launchdarkly/ld-relay-helm/issues/146)) ([e341c29](https://github.com/launchdarkly/ld-relay-helm/commit/e341c2991d3e229689fbffb75c2850f7aa6f7bdf))
* Pass unit tests by fixing golden files ([831fc3f](https://github.com/launchdarkly/ld-relay-helm/commit/831fc3fc654836c8570acc1a4d71f650e024afa8))
* Render the default CPU metric explicitly when no HPA metrics are set ([#145](https://github.com/launchdarkly/ld-relay-helm/issues/145)) ([62f0a91](https://github.com/launchdarkly/ld-relay-helm/commit/62f0a912dc75c4454d8c4ca24eaf525fae228598))
* Use fullname for config resources ([#86](https://github.com/launchdarkly/ld-relay-helm/issues/86)) ([b71f668](https://github.com/launchdarkly/ld-relay-helm/commit/b71f6687e6868fde9dceeaff87d42b7a540bf334)), closes [#85](https://github.com/launchdarkly/ld-relay-helm/issues/85)

## [3.11.4](https://github.com/launchdarkly/ld-relay-helm/compare/3.11.3...3.11.4) (2026-11-19)


### Bug Fixes

* Bump default relay proxy version from 8.23.2 to 8.23.3 ([#156](https://github.com/launchdarkly/ld-relay-helm/issues/156)) ([bda8a41](https://github.com/launchdarkly/ld-relay-helm/commit/bda8a4140192282748f8e7d40a7d7b6abd6285be))
* Bump default relay proxy version from 8.23.3 to 8.23.4 ([#158](https://github.com/launchdarkly/ld-relay-helm/issues/158)) ([194827a](https://github.com/launchdarkly/ld-relay-helm/commit/194827ac3a58814b1f4233d3fcc7f5ce7abb0bb9))

## [3.11.3](https://github.com/launchdarkly/ld-relay-helm/compare/3.11.2...3.11.3) (2026-10-20)


### Bug Fixes

* Bump default relay proxy version from 8.23.1 to 8.23.2 ([#154](https://github.com/launchdarkly/ld-relay-helm/issues/154)) ([3f9011c](https://github.com/launchdarkly/ld-relay-helm/commit/3f9011cff8e0948b50d3e46733ba79c3876802f1))

## [3.11.2](https://github.com/launchdarkly/ld-relay-helm/compare/3.11.1...3.11.2) (2026-09-21)


### Bug Fixes

* Bump default relay proxy version from 8.21.0 to 8.22.0 ([#148](https://github.com/launchdarkly/ld-relay-helm/issues/148)) ([56ea574](https://github.com/launchdarkly/ld-relay-helm/commit/56ea574bb1faae4d44841219677b894d087dbb57))

## [3.11.1](https://github.com/launchdarkly/ld-relay-helm/compare/3.11.0...3.11.1) (2026-09-14)


### Bug Fixes

* hash only the config data in the checksum/config annotation ([#146](https://github.com/launchdarkly/ld-relay-helm/issues/146)) ([e341c29](https://github.com/launchdarkly/ld-relay-helm/commit/e341c2991d3e229689fbffb75c2850f7aa6f7bdf))

## [3.11.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.10.1...3.11.0) (2026-09-08)


### Features

* allow additional HPA metrics via autoscaling.additionalMetrics ([#143](https://github.com/launchdarkly/ld-relay-helm/issues/143)) ([ba05727](https://github.com/launchdarkly/ld-relay-helm/commit/ba0572711d584ee38bf1c6bdfb7efa5cfc033d8d))


### Bug Fixes

* Render the default CPU metric explicitly when no HPA metrics are set ([#145](https://github.com/launchdarkly/ld-relay-helm/issues/145)) ([62f0a91](https://github.com/launchdarkly/ld-relay-helm/commit/62f0a912dc75c4454d8c4ca24eaf525fae228598))

## [3.10.1](https://github.com/launchdarkly/ld-relay-helm/compare/3.10.0...3.10.1) (2026-08-27)


### Bug Fixes

* Bump default relay proxy version from 8.19.1 to 8.21.0 ([#139](https://github.com/launchdarkly/ld-relay-helm/issues/139)) ([348e242](https://github.com/launchdarkly/ld-relay-helm/commit/348e242ce3c20616c8d6e34c044d8629e6d2e4d9))

## [3.10.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.9.0...3.10.0) (2026-07-20)


### Features

* Add support for HTTP routing ([#132](https://github.com/launchdarkly/ld-relay-helm/issues/132)) ([5c79d98](https://github.com/launchdarkly/ld-relay-helm/commit/5c79d9885941a4732482f536b8ebd8a6770b81c7))

## [3.9.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.8.0...3.9.0) (2026-06-29)


### Features

* expose revisionHistoryLimit and dnsConfig/dnsPolicy as configurable values ([#125](https://github.com/launchdarkly/ld-relay-helm/issues/125)) ([c92136f](https://github.com/launchdarkly/ld-relay-helm/commit/c92136f2c59668295cfc6c2e1c2706e5bbe1cf58))

## [3.8.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.7.0...3.8.0) (2026-06-10)


### Features

* expose deployment rolling update strategy and HPA scale behavior as configurable values ([#117](https://github.com/launchdarkly/ld-relay-helm/issues/117)) ([c3b2259](https://github.com/launchdarkly/ld-relay-helm/commit/c3b22590d9f6c9549c04c9d4128c8b214a30c5b1))

## [3.7.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.6.1...3.7.0) (2026-06-02)


### Features

* Add explicit namespace support ([#115](https://github.com/launchdarkly/ld-relay-helm/issues/115)) ([90f85fd](https://github.com/launchdarkly/ld-relay-helm/commit/90f85fdb2d4b16cb9d502a2ff19e1a5719a005a2))

## [3.6.1](https://github.com/launchdarkly/ld-relay-helm/compare/3.6.0...3.6.1) (2026-05-06)


### Bug Fixes

* Bump default relay proxy version from 8.10.5 to 8.19.1 ([#112](https://github.com/launchdarkly/ld-relay-helm/issues/112)) ([b0eecc9](https://github.com/launchdarkly/ld-relay-helm/commit/b0eecc921d6d6ca93950db50ac99564ac152a24b))

## [3.6.0](https://github.com/launchdarkly/ld-relay-helm/compare/3.5.0...3.6.0) (2026-04-06)


### Features

* Add option to configure startup probes in ld-relay container ([#103](https://github.com/launchdarkly/ld-relay-helm/issues/103)) ([cdac3d5](https://github.com/launchdarkly/ld-relay-helm/commit/cdac3d5eb39657e46282c5ab1b7511ac590a7649))


### Bug Fixes

* Default relay.lifecycle to map instead of array ([#99](https://github.com/launchdarkly/ld-relay-helm/issues/99)) ([613e261](https://github.com/launchdarkly/ld-relay-helm/commit/613e261bc287cde3951dc82a6e1afa8a0f9575ce))

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
