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
