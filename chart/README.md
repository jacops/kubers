# Kubers Helm Chart

## Installing the Chart
Add the Flux repo:
```
helm repo add jacops https://charts.jacops.pl
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
| `image.repository` | `docker.io/jacops/kubersctl` | Image repository |
| `image.tag` | `<VERSION>` | Image tag |
| `image.pullPolicy` | `IfNotPresent` | Image pull policy |
| `rbac.enabled` | `true` | If `true`, create and use RBAC resources |
| `injector.serviceAccount.create` | `true` | If `true`, create a new service account |
| `injector.serviceAccount.name` | `kubers` | Service account to be used
| `injector.resources.requests.cpu` | `50m` | CPU resource requests for the Flux deployment |
| `injector.resources.requests.memory` | `64Mi` | Memory resource requests for the Flux deployment | | `injector.resources.limits` | `None` | CPU/memory resource limits for the Flux deployment |
