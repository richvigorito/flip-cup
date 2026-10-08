# One-time bootstrap, applied by a human with local gcloud credentials.
# Creates everything the GitHub Actions deploy needs before it can run:
# APIs, Terraform state bucket, Artifact Registry, service accounts (iam.tf),
# Workload Identity Federation (wif.tf) and a budget alert (budget.tf).

locals {
  apis = [
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "sts.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "billingbudgets.googleapis.com",
  ]
}

data "google_project" "this" {
  project_id = var.project_id
}

resource "google_project_service" "apis" {
  for_each           = toset(local.apis)
  service            = each.value
  disable_on_destroy = false
}

resource "google_storage_bucket" "tfstate" {
  name                        = "${var.project_id}-flipcup-tfstate"
  location                    = var.region
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
  force_destroy               = false

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.apis]
}

resource "google_artifact_registry_repository" "flipcup" {
  repository_id = "flipcup"
  location      = var.region
  format        = "DOCKER"
  description   = "FlipCup container images"

  # Keep the registry inside the 0.5 GB free tier.
  cleanup_policies {
    id     = "keep-recent"
    action = "KEEP"

    most_recent_versions {
      keep_count = 5
    }
  }

  cleanup_policies {
    id     = "delete-old"
    action = "DELETE"

    condition {
      older_than = "604800s"
    }
  }

  depends_on = [google_project_service.apis]
}
