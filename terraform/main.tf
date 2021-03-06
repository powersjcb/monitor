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

resource "google_app_engine_domain_mapping" "domain_mapping" {
  project = var.project_id
  domain_name = "${var.app_name}.${var.domain}"

  ssl_settings {
    ssl_management_type = "AUTOMATIC"
  }
}

resource "google_dns_managed_zone" "public" {
  project      = var.project_id

  dns_name    = "${var.domain}."
  description = var.domain
  name        = "public-zone"
}

// this requires setting up the glcoud provider account email that terraform is using (not default)
resource "google_dns_record_set" "app_cname" {
  project = var.project_id

  managed_zone = google_dns_managed_zone.public.name
  name = "monitor.jacobpowers.me."
  rrdatas = ["ghs.googlehosted.com."]
  ttl = 60
  type = "CNAME"
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

resource "google_secret_manager_secret" "jwt_ec_private_key" {
  provider = google-beta
  secret_id = "${var.app_name}_jwt_ec_private_key"

  replication {
    automatic = true
  }
  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret" "jwt_ec_public_key" {
  provider = google-beta
  secret_id = "${var.app_name}_jwt_ec_public_key"

  replication {
    automatic = true
  }
  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret" "google_client_id" {
  provider = google-beta
  secret_id = "${var.app_name}_google_client_id"

  replication {
    automatic = true
  }
  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret" "google_client_secret" {
  provider = google-beta
  secret_id = "${var.app_name}_google_client_secret"

  replication {
    automatic = true
  }
  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret" "google_client_redirect_url" {
  provider = google-beta
  secret_id = "${var.app_name}_google_client_redirect_url"

  replication {
    automatic = true
  }
  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret" "app_api_key" {
  provider = google-beta
  secret_id = "${var.app_name}_api_key"

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

resource "google_secret_manager_secret_iam_member" "app_engine_jwt_ec_public_key" {
  provider = google-beta
  project = var.project_id

  secret_id = google_secret_manager_secret.jwt_ec_public_key.id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.app_engine.email}"
}

resource "google_secret_manager_secret_iam_member" "app_engine_jwt_ec_private_key" {
  provider = google-beta
  project = var.project_id

  secret_id = google_secret_manager_secret.jwt_ec_private_key.id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.app_engine.email}"
}

resource "google_secret_manager_secret_iam_member" "app_engine_google_client_id" {
  provider = google-beta
  project = var.project_id

  secret_id = google_secret_manager_secret.google_client_id.id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.app_engine.email}"
}

resource "google_secret_manager_secret_iam_member" "app_engine_google_client_secret" {
  provider = google-beta
  project = var.project_id

  secret_id = google_secret_manager_secret.google_client_secret.id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.app_engine.email}"
}

resource "google_secret_manager_secret_iam_member" "app_engine_google_client_redirect_url" {
  provider = google-beta
  project = var.project_id

  secret_id = google_secret_manager_secret.google_client_redirect_url.id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.app_engine.email}"
}

resource "google_secret_manager_secret_iam_member" "app_engine_app_api_key" {
  provider = google-beta
  project = var.project_id

  secret_id = google_secret_manager_secret.app_api_key.id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.app_engine.email}"
}