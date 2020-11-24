variable "region" {
  type        = string
  default     = "westeurope"
  description = "Region: azure region"
}

variable "dns_zone_name" {
  description = "The domain name to create a DNS zone for"
}

variable "subscription_id" {}

variable "tenant_id" {}