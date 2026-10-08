variable "project_id" {
  type        = string
  description = "GCP project ID that hosts FlipCup."
}

variable "region" {
  type        = string
  default     = "us-west1"
  description = "Cloud Run region."
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
  type        = string
  default     = "flipcup.dev"
  description = "Custom domain mapped to the Cloud Run service."
}
