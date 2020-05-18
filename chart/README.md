# Kubers Helm Chart

## Installing the Chart
Add the Flux repo:
```
helm repo add jacops https://jacops.github.io/helm-charts/
```

### Install the chart with the release name `kubers`
1. Create the flux namespace:
   ```
   kubectl create namespace kubers
   ```

2. Install release:
   ```
   helm upgrade -i kubers jacops/kubers \
      --wait --namespace kubers
   ```

## Configuration

| Parameter | Default | Description |
| --- | --- | --- |
| `agent.image.repository` | `docker.io/jacops/kubers-agent` | Image repository |
| `agent.image.tag` | `<VERSION>` | Image tag |
| `agent.provider.name` | `` | If set, the provider annotation doesn't need to be set on pods |
| `agent.provider.aws.region` | `` | AWS region for AWS provider |
| `injector.image.repository` | `docker.io/jacops/kubersd` | Image repository |
| `injector.image.tag` | `<VERSION>` | Image tag |
| `injector.image.pullPolicy` | `IfNotPresent` | Image pull policy |
| `injector.image.log.format` | `standard` | Log output format |
| `injector.image.log.level` | `info` | Log verbosity level. Supported values (in order of detail) are "trace", "debug", "info", "warn", and "err" |
| `injector.serviceAccount.create` | `true` | If `true`, create a new service account |
| `injector.serviceAccount.name` | `kubers` | Service account to be used
| `injector.resources.requests.cpu` | `50m` | CPU resource requests for the Flux deployment |
| `injector.resources.requests.memory` | `64Mi` | Memory resource requests for the Flux deployment |
| `injector.resources.limits` | `None` | CPU/memory resource limits for the Flux deployment |
| `injector.webhook.port` | `8080` | Address to bind listener to |
| `injector.webhook.tls.auto` | `true` | Should the self-signed certs be used for the mutation webhook |
| `injector.webhook.tls.caBundle` | `` | CA bundle for webhook. Should be used in conjunction with `certFile` and `keyFile` |
| `injector.webhook.tls.certFile` | `` | PEM-encoded TLS certificate to serve |
| `injector.webhook.tls.keyFile` | `` | PEM-encoded TLS private key to serve |
| `rbac.enabled` | `true` | If `true`, create and use RBAC resources |
