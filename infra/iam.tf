resource "google_cloud_run_v2_service_iam_member" "public" {
  name     = google_cloud_run_v2_service.flipcup.name
  location = google_cloud_run_v2_service.flipcup.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}
