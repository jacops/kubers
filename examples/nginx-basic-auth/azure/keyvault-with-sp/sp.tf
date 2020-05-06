resource "azuread_application" "nginx_identity" {
  name = join("-", [local.resource_basename, random_id.unique_resource.hex, "nginx-identity", "sp"])
}

resource "azuread_service_principal" "nginx_identity" {
  application_id = azuread_application.nginx_identity.application_id
}

resource "azuread_service_principal_password" "nginx_identity" {
  service_principal_id = azuread_service_principal.nginx_identity.id
  value                = "VT=uSgbTanZhyz@%nL9Hpd+Tfay_MRV#"
  end_date             = "2099-01-01T01:02:03Z"
}

resource "azurerm_key_vault_access_policy" "nginx_identity" {
  key_vault_id = azurerm_key_vault.main.id

  object_id = azuread_service_principal.nginx_identity.object_id
  tenant_id = data.azurerm_client_config.current.tenant_id

  secret_permissions = ["get"]
}
