provider "google" {
  credentials = file("~/jacobpowers-2405446363b3.json")
  project = var.project
  region = var.region
}

resource "random_id" "instance_id" {
  byte_length = 8
}

resource "google_compute_instance" "default" {
  name = "monitor-vm-${random_id.instance_id.hex}"
  machine_type = "f1-micro"
  zone = var.zone

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }
  metadata_startup_script = "echo 'hello world'"
  metadata = {
    ssh-keys = "powersjcb:${file("~/.ssh/jacobpowers_personal_macbook.pub")}"
  }

  network_interface {
    network = "default"
    // default network config to avoid setting up Cloud NAT
    access_config {
    }
  }
}

output "ip" {
  value = google_compute_instance.default.network_interface.0.access_config.0.nat_ip
}

resource "google_compute_instance_group" "api" {
  project = var.project
  name = "${var.name}-instance-group"
  zone = var.zone
  instances = [google_compute_instance.default.self_link]

  lifecycle {
    create_before_destroy = true
  }
  named_port {
    name = "http"
    port = 5000
  }
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