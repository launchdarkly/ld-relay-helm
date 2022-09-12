# Automatic configuration

> **NOTE** Relay Proxy automatic configuration is an Enterprise feature

With [automatic configuration][auto-config], the Relay Proxy automatically detects and supports new connections to updated or new environments whenever changes occur that impact LaunchDarkly SDK keys, mobile keys, or client-side IDs. Configuration is done in two steps.

1. Create a Relay Proxy configuration from the [Relay proxy tab][proxy-tab] of the Account settings page and save its unique key.
2. Provision the Relay Proxy.

    ```yaml
    # auto-config.yaml
    relay:
       environment:
          AUTO_CONFIG_KEY: relay-proxy-key
    ```

    ```shell
    helm install relay --values ./auto-config.yaml launchdarkly-ld-relay/ld-relay
    ```

To learn more about Automatic Configuration in the Relay Proxy, [read here][auto-config].

[auto-config]: https://docs.launchdarkly.com/home/relay-proxy/automatic-configuration
[proxy-tab]: https://app.launchdarkly.com/settings/relay
