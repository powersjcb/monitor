variable "project_id" {
  type = string
  default = "carbide-datum-276117"
}

variable "google_creds" {
  description = "Filesystem location of goole api key"
  type = string
  default = "~/carbide-datum-276117-aa7fb5dcd251.json"
}

variable "region" {
  description = "The region to create the resources in."
  type        = string
  default     = "us-west2"
}

variable "app_name" {
  description = "The name of the deployment."
  type        = string
  default     = "monitor"
}

variable "zone_char" {
  description = "The char name of the zone"
  type        = string
  default     = "a"
}

variable "domain" {
  description = "root domain for application"
  type = string
  default = "jacobpowers.me"
}