# Production stack: Cloud Run service for FlipCup plus the flipcup.dev mapping.
# Applied by GitHub Actions (.github/workflows/deploy-prod.yml). Run
# infra/bootstrap once first.

terraform {
  required_version = ">= 1.6"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  # Bucket is supplied with: terraform init -backend-config="bucket=<name>"
  backend "gcs" {
    prefix = "flipcup/prod"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

variable "project_id" {
  type = string
}

variable "region" {
  type    = string
  default = "us-west1"
}

variable "image" {
  type        = string
  description = "Full Artifact Registry image reference to deploy."
}

variable "runtime_service_account" {
  type        = string
  description = "Email of the flipcup-runtime service account created by bootstrap."
}

variable "domain" {
  type    = string
  default = "flipcup.dev"
}

# Game state is in memory, so the service must stay a single instance and
# scale to zero when idle. Request-based CPU billing keeps idle cost at zero.
resource "google_cloud_run_v2_service" "flipcup" {
  name     = "flipcup"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  deletion_protection = false

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

resource "google_cloud_run_v2_service_iam_member" "public" {
  name     = google_cloud_run_v2_service.flipcup.name
  location = google_cloud_run_v2_service.flipcup.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# The caller must be a verified owner of the domain in Google Search Console.
resource "google_cloud_run_domain_mapping" "apex" {
  name     = var.domain
  location = var.region

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_v2_service.flipcup.name
  }
}

output "service_url" {
  value = google_cloud_run_v2_service.flipcup.uri
}

output "dns_records" {
  description = "Create these records at the flipcup.dev registrar."
  value       = google_cloud_run_domain_mapping.apex.status[0].resource_records
}
