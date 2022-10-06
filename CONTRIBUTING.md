Contributing to the LaunchDarkly Relay Proxy Helm Chart
================================================

Submitting bug reports and feature requests
------------------

The LaunchDarkly SDK team monitors the [issue tracker](https://github.com/launchdarkly/ld-relay-helm/issues) in the Helm chart repository. Bug reports and feature requests specific to this repository should be filed in this issue tracker. The SDK team will respond to all newly filed issues within two business days.

Submitting pull requests
------------------

We encourage pull requests and other contributions from the community. Before submitting pull requests, ensure that all temporary or unintended code is removed. Don't worry about adding reviewers to the pull request; the LaunchDarkly SDK team will add themselves. The SDK team will acknowledge all pull requests within two business days.

Build instructions
------------------

### Prerequisites

Development of this repository requires both `helm` and `go` are installed.

### Testing

The unit tests start with a base set of comparison files referred to as "golden files". These should represent an
identical match to the installed chart if only the default values are used. Starting from this base point allows us to
force additional unit tests only on the properties we expect to change.

Golden files can be updated by running `make update-golden-files`.
Unit tests are run with `make test`.

> **NOTE** The testing structure drew heavy inspiration from [this blog post](https://camunda.com/blog/2022/03/test/).
