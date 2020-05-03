provider "google" {
  credentials = file(var.google_creds)
  region = var.region
}

resource "random_id" "default" {
  byte_length = 8
}

resource "google_app_engine_application" "app" {
  project = var.project_id
  location_id = var.region
}

resource "google_app_engine_domain_mapping" "domain_mapping" {
  domain_name = "jacobpowers.me"

  ssl_settings {
    ssl_management_type = "AUTOMATIC"
  }
}

resource "google_app_engine_firewall_rule" "rule" {
  project = var.project_id
  priority = 1000
  action = "ALLOW"
  source_range = "*"
}

// DATABASE SETUP

resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "primary" {
  // google config
  name = "${var.app_name}-db-primary-${random_id.db_name_suffix.hex}-instance"
  region = var.region
  project = var.project_id

  // database config
  database_version = "POSTGRES_12"

  settings {
    tier = "db-f1-micro"

    disk_autoresize = true
    disk_size = 10

    ip_configuration {
      ipv4_enabled = true
      authorized_networks {
        name = "all-ipv4"
        value = "0.0.0.0/0"
      }
    }
  }
}

resource "random_password" "db_password" {
  length = "32"
  special = false
}

resource "google_sql_user" "app" {
  project = var.project_id
  instance = google_sql_database_instance.primary.name
  name = "app"
  password = random_password.db_password.result
}

resource "google_sql_database" "primary" {
  project = var.project_id
  name = "${var.app_name}-db-primary-${random_id.db_name_suffix.hex}"
  instance = google_sql_database_instance.primary.name
}