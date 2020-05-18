#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PARENT_DIR="$( cd "$( dirname $DIR )" >/dev/null 2>&1 && pwd )"

function tf_value() {
  cd "${PARENT_DIR}"
  printf "%s" $(terraform output $1)
}

function tf_value_base64() {
  cd "${PARENT_DIR}"
  printf "%s" $(terraform output $1) | base64
}

export KEY_VAULT_NAME=$(tf_value "key_vault_name")
export KEY_VAULT_KEY=$(tf_value "key_vault_key")

export ARM_TENANT_ID=$(tf_value_base64 "tenantId")
export ARM_SUBSCRIPTION_ID=$(tf_value_base64 "subscriptionId")
export ARM_CLIENT_ID=$(tf_value_base64 "clientId")
export ARM_CLIENT_SECRET=$(tf_value_base64 "clientSecret")

cd "$DIR"
kubectl kustomize | envsubst
