# Automatic configuration

> **NOTE** Relay Proxy automatic configuration is an Enterprise feature

With [automatic configuration][auto-config], the Relay Proxy automatically detects and supports new connections to updated or new environments whenever changes occur that impact LaunchDarkly SDK keys, mobile keys, or client-side IDs. It is designed for ease of use and minimal maintenance.

Configuration through helm is no different, requiring only two steps.

1. Create a Relay Proxy configuration from the [Relay proxy tab][proxy-tab] of the Account settings page and save its unique key.
2. Configure the Relay Proxy with this key.

    ```shell
    helm install relay --set relay.environment.AUTO_CONFIG_KEY=relay-proxy-key launchdarkly-ld-relay/ld-relay
    ```

    This example seeks to illustrate how minimal the required configuration is by specifying it directly on the command line. Please refer to the [Getting Started](../getting-started.md) documentation for alternative configuration options.

To learn more about Automatic Configuration in the Relay Proxy, [read here][auto-config].

[auto-config]: https://docs.launchdarkly.com/home/relay-proxy/automatic-configuration
[proxy-tab]: https://app.launchdarkly.com/settings/relay