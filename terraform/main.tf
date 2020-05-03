provider "google" {
  credentials = file("~/jacobpowers-2405446363b3.json")
  project = var.project
  region = var.region
}

resource "random_id" "instance_id" {
  byte_length = 8
}

// allow load balancers access to instance
resource "google_compute_firewall" "firewall" {
  name = "monitor-firewall"
  network = "default"

  // load balancer IP ranges
  // https://cloud.google.com/load-balancing/docs/https/#firewall_rules
  source_ranges = ["130.211.0.0/22", "35.191.0.0/16"]

  target_tags = ["monitor-app"]
  source_tags = ["monitor-app"]

  allow {
    protocol = "tcp"
    ports = ["5000"]
  }
}