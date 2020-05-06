resource "azurerm_key_vault" "main" {
  name                = join("-", [local.resource_basename, random_id.unique_resource.hex, "kv"])
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "this" {
  key_vault_id = azurerm_key_vault.main.id

  object_id = data.azurerm_client_config.current.object_id
  tenant_id = data.azurerm_client_config.current.tenant_id

  secret_permissions = ["delete", "get", "list", "set"]
}

resource "azurerm_key_vault_secret" "nginx_htpasswd" {
  name  = "nginx-htpasswd"
  value = "test:$2y$10$i1yZDJmNSJkipWW1JDSp6uIpgt12oCdsM5Q8FZcntOBLerl5dFQ56" #user test:test

  key_vault_id = azurerm_key_vault_access_policy.this.key_vault_id
}
