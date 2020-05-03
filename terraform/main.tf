resource "random_id" "default" {
  byte_length = 8
}

provider "google" {
  credentials = file(var.google_creds)
  region = var.region
}

resource "google_app_engine_application" "app" {
  project = var.project_id
  location_id = var.region
}

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