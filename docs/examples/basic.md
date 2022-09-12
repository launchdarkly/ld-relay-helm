# Minimal example

A successful deployment of the relay proxy requires a single piece of information -- a valid SDK key. This can be specified either directly on the command line,

```shell
helm install relay --set relay.environment.LD_ENV_MyEnvironment=your-sdk-key launchdarkly-ld-relay/ld-relay
```

or by setting the values from a file.

```yaml
# values.yaml
relay:
  environment:
    LD_ENV_MyEnvironment: your-sdk-key
```

```shell
helm install relay --values ./values.yaml launchdarkly-ld-relay/ld-relay
```

This second method is better for long term maintenance, and all subsequent examples in this documentation will favor that approach.

You can specify any valid environment variable directly in the values file to configure the relay proxy instance. To learn more about configuring the Relay Proxy, [read here][proxy-config].

## Using secrets

If you are tracking `values.yaml` in source control, you may not want to store your SDK key in plaintext. This chart provides the option to load environment variables from Kubernetes secrets as well.

1. Create a secret in Kubernetes for storing your SDK key.

    ```shell
    kubectl create secret generic relay --from-literal=sdk-key=your-sdk-key
    ```

2. Provision the Relay Proxy.

    ```yaml
    # values.yaml
    relay:
      secrets:
        - envName: LD_ENV_MyEnvironment
          secretName: relay
          secretKey: sdk-key
    ```

    ```shell
    helm install relay --values ./values.yaml launchdarkly-ld-relay/ld-relay
    ```

[proxy-config]: https://github.com/launchdarkly/ld-relay/blob/v6/docs/configuration.md
