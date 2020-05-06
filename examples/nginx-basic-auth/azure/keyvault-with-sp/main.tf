terraform {
  required_version = ">= 0.12.2"

  backend "local" {}
}

provider "azurerm" {
  version = "~> 2.7"

  features {}
}

provider "random" {
  version = "~> 2.2"
}

locals {
  resource_basename = "akv-k8s"
}

data "azurerm_client_config" "current" {}

resource "random_id" "unique_resource" {
  byte_length = 2
  keepers     = {
    subscription_id = data.azurerm_client_config.current.subscription_id
  }
}

resource "azurerm_resource_group" "main" {
  name     = local.resource_basename
  location = var.location
}
