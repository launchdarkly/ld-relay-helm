# Offline mode

> **Relay Proxy offline mode is an Enterprise feature**. To learn more, read [Offline mode][offline].

Enabling offline mode on the Relay Proxy lets you run the Relay Proxy without ever connecting it to LaunchDarkly. Instead of retrieving flag and segment values from LaunchDarkly's servers, the Relay Proxy gets them from files located on your local host or filesystem.

When using this Helm chart, the offline file needs to exist in a Kubernetes volume which is mounted to the Relay Proxy container. This volume can be created using a [local volume mount][local-volume] or a [ConfigMap][configmap].

To get started, create a Relay Proxy configuration from the [Relay proxy tab][proxy-tab] of the Account settings page in the LaunchDarkly user interface, and save its unique key. Then follow the instructions below to set up the Relay Proxy in offline mode.

## From a local volume

2. A local volume requires a file on the minikube host. Connect to minikube and download a local copy of the flag and segment data using the key from the previous step.

    ```shell
    $ minikube ssh
    $ pwd
    /home/docker
    $ curl https://sdk.launchdarkly.com/relay/latest-all \
      -H "Authorization: rel-EXAMPLE-RELAY-PROXY-CONFIGURATION-KEY" \
      -o EXAMPLE-NAME-OF-OUTPUTTED-FILE.tar.gz
    ```

3. Create a volume and associated volume claim. This allows access to this file within the cluster.

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

4. Now that you have a volume accessible file, configure `values.yaml` to reference this volume claim.

    ```yaml
    # values.yaml
    relay:
      volume:
        # This filename should match the path of the file in the volume used in the
        # below claim.
        offline: EXAMPLE-NAME-OF-OUTPUTTED-FILE.tar.gz
        definition:
          persistentVolumeClaim:
            claimName: offline-volume-claim
    ```

5. Install the Helm chart, referencing your updated values configuration file.

    ```shell
    helm install relay --values ./values.yaml launchdarkly-ld-relay/ld-relay
    ```

## Using a ConfigMap

2. Download a local copy of the flag and segment data using the key from the previous step. Note that unlike with local volumes, this can be done from your local machine.

    ```shell
    $ curl https://sdk.launchdarkly.com/relay/latest-all \
      -H "Authorization: rel-EXAMPLE-RELAY-PROXY-CONFIGURATION-KEY" \
      -o EXAMPLE-NAME-OF-OUTPUTTED-FILE.tar.gz
    ```
3. Create a configmap from this file.

    ```shell
    kubectl create configmap offline-configmap --from-file=EXAMPLE-NAME-OF-OUTPUTTED-FILE.tar.gz
    ```

4. Configure `values.yaml` to reference this configmap.

    ```yaml
    # values.yaml
    relay:
        volume:
            offline: EXAMPLE-NAME-OF-OUTPUTTED-FILE.tar.gz
            definition:
                configMap:
                    name: offline-configmap
    ```

5. Install the Helm chart, referencing your updated values configuration file.

    ```shell
    helm install relay --values ./values.yaml launchdarkly-ld-relay/ld-relay
    ```


Success! Now you should have a working installation of the Relay Proxy, initially configured directly from your pre-downloaded offline file.

[minikube]: https://minikube.sigs.k8s.io/docs/start/
[offline]: https://docs.launchdarkly.com/home/relay-proxy/offline
[proxy-tab]: https://app.launchdarkly.com/settings/relay
[local-volume]: https://kubernetes.io/docs/concepts/storage/volumes/#local
