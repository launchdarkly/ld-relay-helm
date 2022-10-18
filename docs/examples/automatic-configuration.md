# Automatic configuration

> **Relay Proxy automatic configuration is an Enterprise feature**. To learn more, read [Automatic configuration][auto-config].

With automatic configuration, the Relay Proxy automatically detects and supports new connections to updated or new environments whenever changes occur that impact LaunchDarkly SDK keys, mobile keys, or client-side IDs. It is designed for ease of use and minimal maintenance.

Here's how to use automatic configuration through Helm:

1. Create a Relay Proxy configuration from the [Relay proxy tab][proxy-tab] of the Account settings page in the LaunchDarkly user interface, and save its unique key.
2. Configure the Relay Proxy with this key.

    ```shell
    helm install relay --set relay.environment.AUTO_CONFIG_KEY=relay-proxy-key launchdarkly-ld-relay/ld-relay
    ```

    This example illustrates how minimal the required configuration is by specifying it directly on the command line. Refer to the [Getting Started](../getting-started.md) documentation for alternative configuration options.

To learn more about Automatic Configuration in the Relay Proxy, read [Automatic configuration][auto-config].

[auto-config]: https://docs.launchdarkly.com/home/relay-proxy/automatic-configuration
[proxy-tab]: https://app.launchdarkly.com/settings/relay
