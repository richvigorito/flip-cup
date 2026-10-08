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

variable "billing_account_id" {
  type        = string
  description = "Billing account ID (gcloud billing accounts list) for the budget alert."
}

variable "budget_usd" {
  type        = number
  default     = 5
  description = "Monthly budget in USD. Alerts at 50%, 90% and 100% of actual spend."
}
