# Game state is in memory, so the service must stay a single instance and
# scale to zero when idle. Request-based CPU billing keeps idle cost at zero.
resource "google_cloud_run_v2_service" "flipcup" {
  name     = "flipcup"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  deletion_protection = false

  # Service-level scaling is API-defaulted; declare it so plans stay clean.
  scaling {
    manual_instance_count = 0
    min_instance_count    = 0
  }

  template {
    service_account                  = var.runtime_service_account
    timeout                          = "3600s"
    max_instance_request_concurrency = 80

    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    containers {
      image = var.image

      ports {
        container_port = 8080
      }

      resources {
        cpu_idle          = true
        startup_cpu_boost = true

        limits = {
          cpu    = "1"
          memory = "256Mi"
        }
      }

      startup_probe {
        tcp_socket {
          port = 8080
        }
        period_seconds    = 1
        failure_threshold = 10
      }
    }
  }
}
