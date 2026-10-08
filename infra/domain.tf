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
