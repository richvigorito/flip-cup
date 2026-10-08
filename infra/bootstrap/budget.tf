# Alerts by email to billing admins; it does not cap spending.
resource "google_billing_budget" "flipcup" {
  billing_account = var.billing_account_id
  display_name    = "flipcup monthly budget"

  budget_filter {
    # The API stores the project number, so use it to avoid a perpetual diff.
    projects = ["projects/${data.google_project.this.number}"]
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units         = tostring(var.budget_usd)
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }

  threshold_rules {
    threshold_percent = 0.9
  }

  threshold_rules {
    threshold_percent = 1.0
  }

  depends_on = [google_project_service.apis]
}
