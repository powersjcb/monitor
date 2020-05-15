provider "google-beta" {
  credentials = file(var.google_creds)
  region = var.region
}

// comment out if bootstrapping and terraform bucket doesnt exist yet
terraform {
  backend "gcs" {
    bucket = "workspace_terraform_state"
    prefix = "terraform/state"
  }
}

// maintain bucket for terraform backend
resource "google_storage_bucket" "tf_state" {
  project       = var.project_id
  name          = "workspace_terraform_state"
  location      = "us-west2"
  requester_pays = false

  versioning {
    enabled = true
  }
}

// setup application
resource "random_id" "default" {
  byte_length = 8
}

resource "google_app_engine_application" "app" {
  project = var.project_id
  location_id = var.region
}

// todo: configure when domain finishes transfer from amazon -> google
//resource "google_app_engine_domain_mapping" "domain_mapping" {
//  project = var.project_id
//  domain_name = "jacobpowers.me"
//
//  ssl_settings {
//    ssl_management_type = "AUTOMATIC"
//  }
//}

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

// ******
// permit app engine to read secrets
// ******

// note: this must be imported after running "gloud app deploy"
resource "google_service_account" "app_engine" {
  project = var.project_id
  account_id = var.project_id
  display_name = "App Engine default service account"
}

// enables secret manager api
resource "google_project_service" "secretmanager" {
  provider = google-beta
  service  = "secretmanager.googleapis.com"
}

resource "google_secret_manager_secret" "db_connection" {
  provider = google-beta
  secret_id = "${var.app_name}_db_connection"

  replication {
    automatic = true
  }
  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret" "hc_api_key" {
  provider = google-beta
  secret_id = "${var.app_name}_hc_api_key"

  replication {
    automatic = true
  }
  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret_iam_member" "app_engine" {
  provider = google-beta
  project = var.project_id

  secret_id = google_secret_manager_secret.db_connection.id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.app_engine.email}"
}

resource "google_secret_manager_secret_iam_member" "app_engine_hc_api_key" {
  provider = google-beta
  project = var.project_id

  secret_id = google_secret_manager_secret.hc_api_key.id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.app_engine.email}"
}