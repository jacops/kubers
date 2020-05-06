resource "local_file" "kustomization" {
    content  = templatefile("${path.cwd}/templates/kustomization.yaml", {
      tenant_id       = data.azurerm_client_config.current.tenant_id
      subscription_id = data.azurerm_client_config.current.subscription_id
      client_id       = azuread_application.nginx_identity.application_id
      client_secret   = azuread_service_principal_password.nginx_identity.value
    })
    filename = "${path.cwd}/dist/kustomization.yaml"
}

resource "local_file" "patch_deployment" {
    content  = templatefile("${path.cwd}/templates/patches/deployment.yaml", {
      key_vault_name = azurerm_key_vault.main.name
      secret_key     = azurerm_key_vault_secret.nginx_htpasswd.name
    })
    filename = "${path.cwd}/dist/patches/deployment.yaml"
}
