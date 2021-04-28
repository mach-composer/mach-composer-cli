variable "region" {
  type        = string
  default     = "europe-west4"
  description = "Region: gcp region"
}

variable "project_id" {
  type        = string
  description = "The Google shared project id"
}

variable "dns_zone_name" {
  type        = string
  description = "DNS Zone name"
}
