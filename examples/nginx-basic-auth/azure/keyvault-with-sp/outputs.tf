output "tenantId" {
  value = data.azurerm_client_config.current.tenant_id
}

output "subscriptionId" {
  value = data.azurerm_client_config.current.subscription_id
}

output "clientId" {
  value = azuread_application.nginx_identity.application_id
}

output "clientSecret" {
  value = azuread_service_principal_password.nginx_identity.value
  sensitive = true
}

output "key_vault_name" {
  value = azurerm_key_vault.main.name
}

output "key_vault_key" {
  value = azurerm_key_vault_secret.nginx_htpasswd.name
}
