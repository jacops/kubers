# Nginx Basic Auth Example
This is a simple demonstration of how `kubers` can be leveraged to protect a website served by `nginx`.
If you are using the provided Terraform scripts, the credentials to access example website are: `test:test`

## Prerequisites
* Kubernetes cluster

## Integration
Instructions on how to integrate this example with some of the secret stores.

> Please install `kubers` first, before you run the examples.

### Azure KeyVault

#### Service Principal method
This is a default integration for this example.

```
KEY_VAULT_NAME=xxx \
KEY_VAULT_KEY=xxx \
ARM_TENANT_ID=xxx \
ARM_SUBSCRIPTION_ID=xxx \
ARM_CLIENT_ID=xxx \
ARM_CLIENT_SECRET=xxx \
kubectl kustomize azure/keyvault-with-sp/deploy | envsubst | kubectl apply -f -
```

If you don't have required infrastructure provisioned, you can use the terraform scripts from `azure/keyvault-with-sp` directory.

Running the command below will automatically substitute environmental variables in yaml files:

```
./azure/keyvault-with-sp/deploy/kustomize.sh | kubectl apply -f -
```

### AWS Secret Manager

#### API Access Keys method
This is a default integration for this example.

```
NGINX_HTPASSWD_SECRET_MANAGER_KEY=xxx \
AWS_ACCESS_KEY_ID=xxx \
AWS_SECRET_ACCESS_KEY=xxx \
kubectl kustomize aws/secretmanager-with-user/deploy | envsubst | kubectl apply -f -
```

If you don't have required infrastructure provisioned, you can use the terraform scripts from `aws/secretmanager-with-user` directory.

Running the command below will automatically substitute environmental variables in yaml files:

```
./aws/secretmanager-with-user/deploy/kustomize.sh | kubectl apply -f -
```
