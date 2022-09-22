# Getting started

The relay proxy container supports source configuration from both files and environment variables. These methods can be used individually or together. Each are fully supported with this helm chart. Let's look at examples showcasing each method.

A successful Relay Proxy deployment minimally requires a valid SDK key. Providing this value to helm can be done directly on the command line.

```shell
helm install relay --set relay.environment.LD_ENV_MyEnvironment=your-sdk-key launchdarkly-ld-relay/ld-relay
```

This method is not well-suited for long term maintenance and quickly grows unforgivably verbose. As an alternative, helm values can be provided through a values file:

```shell
helm install relay --values ./values.yaml launchdarkly-ld-relay/ld-relay
```

The rest of the documentation will assume this is the installation method of choice. Unless otherwise needed, we will omit the repeated `install` command.

## Environment variables

Using environment variables to configure the relay proxy is the easiest way to get started with a deployment. A minimal example might look like the following.

```yaml
# values.yaml
relay:
  environment:
    LD_ENV_MyEnvironment: your-sdk-key
```
Any name/value pair defined here will be set into the relay container. See the [relay proxy's guide on configuration][proxy-config] for additional details.

Setting environment variables directly in the file is convenient, but not always ideal. If you are tracking `values.yaml` in source control, you may not wish to store your SDK key in plaintext.

### Sharing secrets

Instead of embedding sensitive information in our `values.yaml` file, we can store these values independently into the Kubernetes cluster using [secrets].

Because Secrets can be created independently of the Pods that use them, there is less risk of the Secret (and its data) being exposed during the workflow of creating, viewing, and editing Pods

Secrets are created independently of the helm chart installation. To store our SDK key, we can create a secret like:

```shell
kubectl create secret generic relay --from-literal=sdk-key=your-sdk-key
```

This secret name and key can be referenced in our `values.yaml` file.

```yaml
# values.yaml
relay:
  secrets:
    - envName: LD_ENV_MyEnvironment
      secretName: relay
      secretKey: sdk-key
```

Both plaintext variables and secrets can be used together to control the proxy.

```yaml
# values.yaml
relay:
  environment:
    USE_EVENTS: true
  secrets:
    - envName: LD_ENV_MyEnvironment
      secretName: relay
      secretKey: sdk-key
```

## Configuration file

As shown in the [relay proxy's guide][proxy-config], configuration can be controlled using an `ini` style configuration file. Using an existing file with this helm chart requires the use of [volumes].

Our example below makes use of [minikube] and a [local volume mount][local-volume] though any kubernetes volume type can be substituted.

1. A local volume requires a file on the minikube host. We connect to minikube and create our minimal configuration file.

    ```shell
    $ minikube ssh
    $ pwd
    /home/docker
    $ vim relay-proxy-config.conf
    ```

   ```ini
   # relay-proxy-config.conf
   [Environment "MyProduction"]
   sdkKey="your-sdk-key"
   ```

2. Next, we create a volume and associated volume claim. This allows access to this file within the cluster.

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

3. Now that we have a volume accessible file, we can configure `values.yaml` to reference this volume claim.

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

Using this claim, the local volume will be mounted into the running docker image and used as the basis for configuration. As noted in the [relay proxy documentation][proxy-config], any conflicting environment variables will take precedence over the configuration file.

[proxy-config]: https://github.com/launchdarkly/ld-relay/blob/v6/docs/configuration.md
[minikube]: https://minikube.sigs.k8s.io/docs/start/
[volumes]: https://kubernetes.io/docs/concepts/storage/volumes/
[local-volume]: https://kubernetes.io/docs/concepts/storage/volumes/#local
[secrets]: https://kubernetes.io/docs/concepts/configuration/secret/