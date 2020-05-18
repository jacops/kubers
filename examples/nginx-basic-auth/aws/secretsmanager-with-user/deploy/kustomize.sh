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

export NGINX_HTPASSWD_SECRETS_MANAGER_KEY=$(tf_value "nginx_htpasswd_secret_name")
export AWS_ACCESS_KEY_ID=$(tf_value_base64 "nginx_identity_access_key_id")
export AWS_SECRET_ACCESS_KEY=$(tf_value_base64 "nginx_identity_secret_access_key")

cd "$DIR"
kubectl kustomize | envsubst
