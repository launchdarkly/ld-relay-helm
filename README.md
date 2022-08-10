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

## Configuration

This chart can be customized by overriding the configuration options shown in [the values file](https://github.com/launchdarkly/ld-relay-helm/blob/main/values.yaml).

The `relay.environment` configuration option controls the configuration of the underlying relay proxy instance. See the [relay proxy's guide on configuration](https://github.com/launchdarkly/ld-relay/blob/v6/docs/configuration.md#file-section-environment-name) for a list of valid environment variable names.

```shell
# Minimal example specifying only the environment
helm install --set LD_ENV_Production=your-sdk-key ld-relay launchdarkly-ld-relay/ld-relay
```

## Learn more

Check out our [documentation](https://docs.launchdarkly.com) for in-depth instructions on configuring and using LaunchDarkly. You can also head straight to the [complete reference guide for the relay proxy](https://docs.launchdarkly.com/home/relay-proxy).

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
