# Offline mode

> **NOTE** Relay Proxy offline mode is an Enterprise feature

Enabling [offline mode][offline] on the Relay Proxy lets you run the Relay Proxy without ever connecting it to LaunchDarkly. Instead of retrieving flag and segment values from LaunchDarkly's servers, the Relay Proxy gets them from files located on your local host or filesystem. When using this helm chart, this file will need to exist in a Kubernetes volume which will be mounted to the relay proxy container. Our example below makes use of [minikube] and a [local volume mount][local-volume].

1. Create a Relay Proxy configuration from the [Relay proxy tab][proxy-tab] of the Account settings page and save its unique key.
2. Connect to the minikube host and download a local copy of the flag and segment data using the key from the previous step.

    ```shell
    minikube ssh
    curl https://sdk.launchdarkly.com/relay/latest-all \
      -H "Authorization: rel-EXAMPLE-RELAY-PROXY-CONFIGURATION-KEY" \
      -o EXAMPLE-NAME-OF-OUTPUTTED-FILE.tar.gz
    ```

3. Create a volume and associated claim to allow this file to be accessed within the cluster.

    ```yaml
    # offline-volume.yaml
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: offline-volume
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

    # offline-claim.yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: offline-volume-claim
    spec:
      storageClassName: manual
      accessModes:
        - ReadOnlyMany
      resources:
        requests:
          storage: 1Gi
    ```

    ```shell
    kubectl apply -f offline-volume.yaml
    kubectl apply -f offline-claim.yaml
    ```

4. Provision the Relay Proxy.

    ```yaml
    # offline.yaml
    relay:
      offline:
        enabled: true
        # This filename should match the path of the file in the volume used in the
        # below claim.
        filename: relay-file.tar.gz
        volume:
          persistentVolumeClaim:
            claimName: offline-volume-claim
    ```

    ```shell
    helm install relay --values ./offline.yaml launchdarkly-ld-relay/ld-relay
    ```

To learn more about Offline mode in the Relay Proxy, [read here][offline].

[minikube]: https://minikube.sigs.k8s.io/docs/start/
[offline]: https://docs.launchdarkly.com/home/relay-proxy/offline
[proxy-tab]: https://app.launchdarkly.com/settings/relay
[local-volume]: https://kubernetes.io/docs/concepts/storage/volumes/#local