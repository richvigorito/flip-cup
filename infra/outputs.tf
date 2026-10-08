output "service_url" {
  description = "Default run.app URL of the service."
  value       = google_cloud_run_v2_service.flipcup.uri
}

output "dns_records" {
  description = "Create these records at the flipcup.dev registrar."
  value       = google_cloud_run_domain_mapping.apex.status[0].resource_records
}
