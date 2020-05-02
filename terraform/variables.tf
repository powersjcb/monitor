variable "project" {
  description = "The project ID to create the resources in."
  type        = string
  default     = "jacobpowers-200819"
}

variable "region" {
  description = "The region to create the resources in."
  type        = string
  default     = "us-west1"
}

variable "name" {
  description = "The name of the deployment."
  type        = string
  default     = "monitor"
}

variable "zone" {
  description = "The name of the zone"
  type        = string
  default     = "us-west1-a"
}