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
