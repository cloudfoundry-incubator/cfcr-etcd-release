resource "google_compute_firewall" "jumpbox-to-etcd" {
  name    = "${var.env_id}-jumpbox-to-etcd"
  network = "${google_compute_network.bbl-network.name}"

  source_tags = ["${var.env_id}-jumpbox"]

  allow {
    ports    = ["2379"]
    protocol = "tcp"
  }

  target_tags = ["${var.env_id}-internal"]
}
