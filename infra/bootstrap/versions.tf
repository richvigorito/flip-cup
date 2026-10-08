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

  # The Budgets API needs a quota project when called with user credentials.
  user_project_override = true
  billing_project       = var.project_id
}
