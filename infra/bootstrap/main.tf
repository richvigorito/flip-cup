# One-time bootstrap, applied by a human with local gcloud credentials.
# Creates everything the GitHub Actions deploy needs before it can run:
# APIs, Terraform state bucket, Artifact Registry, runtime/deployer service
# accounts, and Workload Identity Federation for GitHub.

terraform {
  required_version = ">= 1.6"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

variable "project_id" {
  type        = string
  description = "GCP project ID that hosts FlipCup."
}

variable "region" {
  type        = string
  default     = "us-west1"
  description = "Free-tier Cloud Run region closest to Seattle."
}

variable "github_repository" {
  type        = string
  default     = "richvigorito/flip-cup"
  description = "owner/repo allowed to deploy through Workload Identity Federation."
}

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

resource "google_service_account" "runtime" {
  account_id   = "flipcup-runtime"
  display_name = "FlipCup Cloud Run runtime"

  depends_on = [google_project_service.apis]
}

resource "google_service_account" "deployer" {
  account_id   = "flipcup-deployer"
  display_name = "FlipCup GitHub Actions deployer"

  depends_on = [google_project_service.apis]
}

resource "google_project_iam_member" "deployer_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.deployer.email}"
}

resource "google_artifact_registry_repository_iam_member" "deployer_push" {
  repository = google_artifact_registry_repository.flipcup.name
  location   = var.region
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${google_service_account.deployer.email}"
}

resource "google_storage_bucket_iam_member" "deployer_state" {
  bucket = google_storage_bucket.tfstate.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.deployer.email}"
}

resource "google_service_account_iam_member" "deployer_acts_as_runtime" {
  service_account_id = google_service_account.runtime.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.deployer.email}"
}

resource "google_iam_workload_identity_pool" "github" {
  workload_identity_pool_id = "github"
  display_name              = "GitHub Actions"

  depends_on = [google_project_service.apis]
}

resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github.workload_identity_pool_id
  workload_identity_pool_provider_id = "github"
  display_name                       = "GitHub OIDC"

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
    "attribute.ref"        = "assertion.ref"
  }

  # Only this repo, and only its main branch, can mint credentials.
  attribute_condition = "assertion.repository == '${var.github_repository}' && assertion.ref == 'refs/heads/main'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_service_account_iam_member" "github_impersonates_deployer" {
  service_account_id = google_service_account.deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_repository}"
}

output "tf_state_bucket" {
  description = "Set as the TF_STATE_BUCKET GitHub variable."
  value       = google_storage_bucket.tfstate.name
}

output "runtime_service_account" {
  description = "Set as the GCP_RUNTIME_SA GitHub variable."
  value       = google_service_account.runtime.email
}

output "deployer_service_account" {
  description = "Set as the GCP_DEPLOY_SA GitHub variable. Also add it as an owner of flipcup.dev in Search Console."
  value       = google_service_account.deployer.email
}

output "workload_identity_provider" {
  description = "Set as the GCP_WIF_PROVIDER GitHub variable."
  value       = google_iam_workload_identity_pool_provider.github.name
}
