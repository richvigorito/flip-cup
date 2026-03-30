variable "datacenter" {
  type    = string
  default = "rv_homelab"
}

variable "image" {
  type = string
}

variable "hostname" {
  type = string
}

variable "node_name" {
  type    = string
  default = "rv-hstack-node1-pi4"
}

variable "container_port" {
  type    = number
  default = 8080
}

variable "count" {
  type    = number
  default = 1
}

variable "cpu" {
  type    = number
  default = 500
}

variable "memory" {
  type    = number
  default = 512
}

job "flipcup" {
  datacenters = [var.datacenter]
  type        = "service"

  group "app" {
    count = var.count

    constraint {
      attribute = "${node.unique.name}"
      operator  = "="
      value     = var.node_name
    }

    network {
      port "http" {}
    }

    service {
      name     = "flipcup"
      provider = "consul"
      port     = "http"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.flipcup.rule=Host(`${var.hostname}`)",
        "traefik.http.routers.flipcup.entrypoints=web",
      ]

      check {
        type     = "http"
        path     = "/quizzes"
        interval = "15s"
        timeout  = "2s"
      }
    }

    update {
      max_parallel      = 1
      canary            = 0
      auto_revert       = true
      health_check      = "checks"
      min_healthy_time  = "20s"
      healthy_deadline  = "2m"
      progress_deadline = "5m"
    }

    restart {
      attempts = 2
      interval = "30m"
      delay    = "15s"
      mode     = "delay"
    }

    task "app" {
      driver = "docker"

      vault {
        role         = "flipcup-staging"
        env          = false
        disable_file = true
        change_mode  = "restart"
      }

      config {
        image      = var.image
        force_pull = false
        ports      = ["http"]
        network_mode = "host"
      }

      template {
        data = <<EOF
PORT="{{ env "NOMAD_PORT_http" }}"
{{ with secret "secret/data/flipcup/staging" -}}
GAME_CLEANUP_INTERVAL="{{ .Data.data.cleanup_interval }}"
GAME_STALE_AFTER="{{ .Data.data.stale_after }}"
{{- end }}
EOF
        destination = "secrets/flipcup.env"
        env         = true
      }

      resources {
        cpu    = var.cpu
        memory = var.memory
      }
    }
  }
}
