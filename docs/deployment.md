# Deployment guide

This document explains the three environment stories in the repo:

- local development with Docker Compose
- staging on the homelab Pi cluster
- production on Google Cloud Run (flipcup.dev)

The key design choice is that local development keeps its own dev-focused Dockerfiles, while staging and production share the repo-root deployment `Dockerfile`.

## Environment contract

| Environment | Packaging path | Frontend websocket config | Backend port |
| --- | --- | --- | --- |
| Local dev | `docker-compose.yml` with `game-server/Dockerfile` and `ui/Dockerfile` dev targets | `VITE_WS_URL=localhost:8080` for the frontend dev server | `PORT=8080` |
| Staging | repo-root `Dockerfile`, built by the self-hosted runner on the staging Pi | leave `VITE_WS_URL` unset so the browser host is used at runtime | `PORT` comes from Nomad |
| Production | repo-root `Dockerfile`, deployed to Cloud Run by Terraform (`infra/`) from GitHub Actions | leave `VITE_WS_URL` unset so the browser host is used at runtime | `PORT=8080` |

The frontend runtime host logic lives in `ui/src/lib/utils/config.ts`, which is why staging and Cloud Run can rely on the current host instead of baking in a hardcoded deployment URL.

## Local development

For the quickest full-stack loop:

```bash
docker-compose up -d
```

Compose uses:

- `game-server/Dockerfile` with target `dev`
- `ui/Dockerfile` with target `dev`

That means the repo-root deployment `Dockerfile` does **not** interfere with the local developer workflow.

You can also run each app directly:

```bash
cd game-server && go run cmd/flipcup/main.go
cd ui && npm install && npm run dev
```

## Staging on Nomad / Consul / Traefik / Vault

The staging job definition lives at:

- `deploy/nomad/flipcup.nomad.hcl`

The staging deploy workflow lives at:

- `.github/workflows/deploy-staging.yml`

### How staging deployment works

1. `Deploy staging` is triggered in GitHub Actions
2. GitHub schedules the deploy job onto the self-hosted runner labeled `self-hosted`, `Linux`, `ARM64`, `staging`
3. that runner is installed on the homelab staging node
4. the runner checks out the repo, builds the repo-root `Dockerfile`, and runs `nomad job run`
5. Nomad schedules the `flipcup` job on the pinned node
6. the task reads its runtime config from Vault and starts behind Traefik

Important detail: GitHub does **not** push into the home network directly. The self-hosted runner inside the homelab polls GitHub for jobs and executes them locally.

### Key assumptions

- staging runs one allocation because the app still keeps game state in memory
- the job is pinned to one Nomad node so the locally built Docker image is available to the Docker driver
- Traefik routes the chosen hostname to the Nomad service
- the app is served from the same host for HTTP and WebSocket traffic

### Required GitHub repository variables

- `STAGING_HOSTNAME`
- optional `STAGING_NOMAD_DATACENTER`
- optional `STAGING_NOMAD_NODE`

Current defaults used by the workflow:

- `STAGING_NOMAD_DATACENTER=rv_homelab`
- `STAGING_NOMAD_NODE=rv-hstack-node1-pi4`

### Runner requirements

The staging runner should:

- run inside the homelab
- be registered on the repo with labels `self-hosted`, `Linux`, `ARM64`, `staging`
- have working access to Docker and Nomad
- ideally live on the same node named by `STAGING_NOMAD_NODE`

That last point matters because the workflow builds the Docker image locally and then runs Nomad on the same host. If you move the runner elsewhere without introducing a shared image registry, Nomad may schedule onto a node that does not have the built image.

### Vault-backed runtime config

The staging task uses Nomad's native Vault integration.

Current contract:

- auth mount: `auth/jwt-nomad`
- role: `flipcup-staging`
- secret read path in templates: `secret/data/flipcup/staging`
- logical write path for CLI usage: `secret/flipcup/staging`

Seed the secret from a machine that already has Vault access:

```bash
vault kv put secret/flipcup/staging cleanup_interval=30m stale_after=1h
```

Those values are rendered into:

- `GAME_CLEANUP_INTERVAL`
- `GAME_STALE_AFTER`

### Staging verification checklist

After a deploy:

```bash
nomad status flipcup
nomad job allocs flipcup
nomad alloc status <alloc-id>
nomad alloc logs <alloc-id>
curl http://flipcup.homelab/quizzes
```

## Production on Google Cloud Run

Production is `https://flipcup.dev`, served by a Cloud Run service in `us-west1` and managed entirely by Terraform. Every push to `main` runs `.github/workflows/deploy-prod.yml`: validate (Go + Playwright), build the root `Dockerfile`, push to Artifact Registry, then `terraform apply` in `infra/`.

### Runtime shape

- `min-instances=0`: scales to zero when idle, so it costs ~$0 at low traffic (expect a 1–3 second cold start for the first visitor)
- `max-instances=1`: game state is in memory, so two instances would split players across lobbies
- request timeout 3600s: Cloud Run closes WebSockets after 60 minutes, so the client must reconnect and rejoin
- request-based CPU billing, 256Mi memory
- idle lobbies are lost when the instance scales to zero; live games keep it alive through their open WebSockets

### Terraform layout

- `infra/bootstrap/` — applied **once by hand**: APIs, state bucket, Artifact Registry, runtime and deployer service accounts, Workload Identity Federation (only `main` of this repo can deploy; no JSON keys)
- `infra/` — applied by CI: Cloud Run service, public invoker, `flipcup.dev` domain mapping, optional budget alert

### One-time setup

1. Create a GCP project with billing enabled, then:

   ```bash
   gcloud auth application-default login
   cd infra/bootstrap
   terraform init
   terraform apply -var project_id=<your-project-id> -var billing_account_id=<ACCOUNT_ID>
   ```

   The bootstrap state is local and gitignored; keep it somewhere safe or migrate it to the state bucket.

2. In GitHub (Settings → Secrets and variables → Actions → Variables), set from the bootstrap outputs:

   | Variable | Value |
   | --- | --- |
   | `GCP_PROJECT_ID` | project ID |
   | `GCP_REGION` | optional, defaults to `us-west1` |
   | `TF_STATE_BUCKET` | `tf_state_bucket` output |
   | `GCP_RUNTIME_SA` | `runtime_service_account` output |
   | `GCP_DEPLOY_SA` | `deployer_service_account` output |
   | `GCP_WIF_PROVIDER` | `workload_identity_provider` output |

   Also create a GitHub environment named `production`.

3. Verify `flipcup.dev` in [Google Search Console](https://search.google.com/search-console) and add the deployer service account email as an **owner**. Cloud Run domain mapping requires the caller to own the domain.

4. Push to `main` (or run the workflow manually). When it finishes, read the `dns_records` Terraform output (in the job log) and create those records at the flipcup.dev registrar. The managed certificate can take up to an hour after DNS propagates.

5. After a few stable days on Cloud Run, retire Fly: `fly apps destroy flipcup`, then delete `fly.toml` and `Dockerfile.fly`.

### Budget alert

The bootstrap also creates a monthly budget (default $5, alerts at 50%, 90% and 100%) that emails billing admins. It does not cap spending. Pass your billing account when applying the bootstrap:

```bash
terraform apply -var project_id=<id> -var billing_account_id=<ACCOUNT_ID> [-var budget_usd=10]
```

## CI and deployment gates

Validation is defined in:

- `.github/workflows/ci.yml`
- `.github/workflows/validate.yml`

That validation is reused by staging deploys so changes get the same Go/UI/Playwright coverage before rollout.

## Protecting `main`

At minimum, keep `main` protected by rules that:

- require pull requests
- block direct pushes
- require the `tests` check from the `CI` workflow

That keeps staging deploys tied to code that has already passed the normal validation path.
