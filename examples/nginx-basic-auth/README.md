# Nginx Basic Auth Example
This is a simple demonstration of how `kubers` can be leveraged to protect a website served by `nginx`.

## Prerequisites
* Kubernetes cluster

## Integration
Instructions on how to integrate this example with some of the secret stores.

> Please install `kubers` first, before you run the examples.

### Azure KeyVault

#### Service Principal method
This is a default integration for this example.

##### No prior infrastructure
Before deploying the example, make sure you have provisioned infrastructure from `./azure/keyvault-with-sp`
This will also create `kustomize` patches for easier example deployment.

```
kubectl apply -k azure/keyvault-with-sp/dist`
```

##### Existing KeyVault and service principal
Before deploying the example, make sure you have:
* changed the annotation in `deploy/deployment.yaml` to match your existing KeyVault name.
* created a secret called `nginx-sp-creds` in the deployment namespace.

```
kubectl apply -k deploy/
```
