# Full-featured example

In this example, we are going to show a more full-featured example which leverages

* Multi-environment support
* Redis as a persistence store
* Exporting Prometheus metrics

We can use helm to deploy a redis instance into our cluster.

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis
```

Then we can apply the following values when provisioning the Relay Proxy:

```yaml
# values.yaml
service:
   type: ClusterIP
   ports:
      # A port named api is always required
      - port: 8030
        targetPort: 8030
        protocol: TCP
        name: api
      # Prometheus support requires exposing additional ports
      - port: 8031
        targetPort: 8031
        protocol: TCP
        name: prometheus

relay:
   environment:
      LD_ENV_Production: your-production-sdk-key
      LD_PREFIX_Production: ld-production
      LD_ENV_Staging: your-staging-sdk-key
      LD_PREFIX_Staging: ld-staging
      USE_REDIS: true
      REDIS_HOST: redis-master.default.svc.cluster.local
      USE_PROMETHEUS: true
   secrets:
      # The bitnami chart creates this secret for us automatically.
      - envName: REDIS_PASSWORD
        secretName: redis
        secretKey: redis-password
```

```shell
helm install relay --values ./values.yaml launchdarkly-ld-relay/ld-relay
```
