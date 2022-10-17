# Getting started

The Relay Proxy container supports source configuration from both files and environment variables. You can use these methods individually or together. Each is fully supported with this Helm chart.

A successful Relay Proxy deployment minimally requires a valid SDK key. You can provide this value to Helm directly on the command line:

```shell
helm install relay --set relay.environment.LD_ENV_YourEnvironment=your-sdk-key launchdarkly-ld-relay/ld-relay
```

This method is not well-suited for long term maintenance and quickly grows unforgivably verbose. As an alternative, you can also provide Helm values through a values file:

```shell
helm install relay --values ./values.yaml launchdarkly-ld-relay/ld-relay
```

## Environment variables

Using environment variables to configure the Relay Proxy is the easiest way to get started with a deployment.

Here's how:

```yaml
# values.yaml
relay:
  environment:
    LD_ENV_YourEnvironment: your-sdk-key
```
Any name/value pair defined here will be set into the Relay container. To learn more, read the [Relay Proxy configuration guide][proxy-config] for additional details.

Setting environment variables directly in the file is convenient, but not always ideal. If you are tracking `values.yaml` in source control, you may not wish to store your SDK key in plaintext.

### Sharing secrets

Instead of embedding sensitive information in your `values.yaml` file, you can store these values independently into the Kubernetes cluster using [Secrets].

Because Secrets can be created independently of the Pods that use them, there is less risk of the Secret and its data being exposed during the workflow of creating, viewing, and editing Pods

You must create a Secret independently of the Helm chart installation. To store your SDK key, you can create a secret using this command:

```shell
kubectl create secret generic relay --from-literal=sdk-key=your-sdk-key
```

Then, you can reference this secret name and key in your `values.yaml` file:

```yaml
# values.yaml
relay:
  secrets:
    - envName: LD_ENV_YourEnvironment
      secretName: relay
      secretKey: sdk-key
```

You can use both plaintext variables and secrets together to control the Relay Proxy. Here's an example:

```yaml
# values.yaml
relay:
  environment:
    USE_EVENTS: true
  secrets:
    - envName: LD_ENV_YourEnvironment
      secretName: relay
      secretKey: sdk-key
```

## Configuration file

As shown in the [Relay Proxy configuration guide][proxy-config], you can control configuration using an `ini` style configuration file. Using an existing file with this Helm chart requires the use of [volumes].

The example below uses [minikube] and a [local volume mount][local-volume]. You can subsitute any Kubernetes volume type.

Here's how to control Relay Proxy configuration using a minikube and a local volume mount:

1. A local volume requires a file on the minikube host. Connect to minikube and create a minimal configuration file.

    ```shell
    $ minikube ssh
    $ pwd
    /home/docker
    $ vim relay-proxy-config.conf
    ```

   ```ini
   # relay-proxy-config.conf
   [Environment "YourProduction"]
   sdkKey="your-sdk-key"
   ```

2. Create a volume and associated volume claim. This allows access to this file within the cluster.

    ```yaml
    # relay-config-file-volume.yaml
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: relay-config-file-volume
      labels:
        type: local
    spec:
      storageClassName: manual
      capacity:
        storage: 1Gi
      accessModes:
        - ReadOnlyMany
      hostPath:
        path: "/home/docker/"

    # relay-config-file-claim.yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: relay-config-file-volume-claim
    spec:
      storageClassName: manual
      accessModes:
        - ReadOnlyMany
      resources:
        requests:
          storage: 1Gi
    ```

    ```shell
    kubectl apply -f relay-config-file-volume.yaml
    kubectl apply -f relay-config-file-claim.yaml
    ```

3. Now that you have a volume accessible file, configure `values.yaml` to reference this volume claim.

    ```yaml
    # values.yaml
    relay:
      volume:
        # This filename should match the path of the file in the volume used in the
        # below claim.
        config: relay-proxy-config.conf
        definition:
          persistentVolumeClaim:
            claimName: relay-config-file-volume-claim
    ```

Using this claim, the local volume is mounted into the running docker image and used as the basis for configuration. Any conflicting environment variables will take precedence over the configuration file. To learn more, read [Relay Proxy configuration][proxy-config].

[proxy-config]: https://github.com/launchdarkly/ld-relay/blob/v6/docs/configuration.md
[minikube]: https://minikube.sigs.k8s.io/docs/start/
[volumes]: https://kubernetes.io/docs/concepts/storage/volumes/
[local-volume]: https://kubernetes.io/docs/concepts/storage/volumes/#local
[secrets]: https://kubernetes.io/docs/concepts/configuration/secret/
