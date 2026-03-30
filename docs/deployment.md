# Deployment guide

## Environment contract

FlipCup now uses one production-oriented container build at the repository root (`Dockerfile`) and keeps environment differences in runtime configuration instead of separate packaging flows.

| Environment | Packaging path | Frontend websocket config | Backend port |
| --- | --- | --- | --- |
| Local dev | `docker-compose.yml` with service-specific dev Dockerfiles | `VITE_WS_URL=localhost:8080` with runtime fallback to the current LAN host when opened from another device | `PORT=8080` |
| Staging | Root `Dockerfile`, built on the staging Nomad host | Leave `VITE_WS_URL` unset so the UI uses the current host | `PORT=8080` |
| Production | Root `Dockerfile`, pushed to Artifact Registry and deployed to Cloud Run | Leave `VITE_WS_URL` unset so the UI uses the current host | `PORT=8080` |

The frontend local override now lives in `ui/.env.development`. That keeps `localhost:8080` out of production Vite builds while still letting remote devices on the same LAN connect through the browser's current host.

## Local development

Use Docker Compose for the quickest full-stack loop:

```bash
docker-compose up -d
```

You can also work in each app directly:

```bash
cd game-server && go run cmd/flipcup/main.go
cd ui && npm install && npm run dev
```

## Staging on Nomad / Consul / Traefik

The staging job definition lives at `deploy/nomad/flipcup.nomad.hcl`.

Key assumptions:

- staging runs a single FlipCup allocation until the app has shared state
- the allocation is pinned to one Nomad node so the locally built image is available to the Docker driver
- Traefik routes the `.homelab` host alias over the `web` entrypoint to the Nomad service registered in Consul
- the app is exposed on `/` and `/ws` from the same hostname

The `Deploy staging` GitHub Actions workflow does the following after validation passes:

1. runs on a self-hosted GitHub Actions runner inside the homelab
2. builds the root `Dockerfile` locally on that runner
3. verifies Docker and the Nomad CLI are available and can reach the cluster
4. runs `nomad job run` with the staging hostname and node variables

Required GitHub configuration for the workflow:

- repository variable `STAGING_HOSTNAME`
- optional variables `STAGING_NOMAD_DATACENTER`, `STAGING_NOMAD_NODE`
- a self-hosted runner labeled `self-hosted`, `Linux`, `ARM64`, and `staging`

Defaults in this repository assume the current Raspberry Pi cluster shape:

- `STAGING_NOMAD_DATACENTER=rv_homelab`
- `STAGING_NOMAD_NODE=rv-hstack-node1-pi4`

The self-hosted runner should live on the same node named by `STAGING_NOMAD_NODE`. The staging job still uses a locally built Docker image and is pinned to one Nomad node, so building and deploying from a different machine would leave the target node without the image unless you switch to a shared image registry.

Override those variables only if the cluster naming changes or you want the job pinned elsewhere.

FlipCup staging runtime settings are now sourced from Vault instead of Nomad Variables.

The Nomad task authenticates to Vault with:

- auth mount `auth/jwt-nomad`
- role `flipcup-staging`
- KV v2 secret read path `secret/data/flipcup/staging`

Seed the secret from a machine that already has Vault access:

```bash
vault kv put secret/flipcup/staging cleanup_interval=30m stale_after=1h
```

The current job reads:

- `cleanup_interval` → `GAME_CLEANUP_INTERVAL`
- `stale_after` → `GAME_STALE_AFTER`

`PORT` still comes from Nomad's allocated service port.

The template stanza reads the KV v2 API path `secret/data/flipcup/staging`, while the CLI command above writes to the logical path `secret/flipcup/staging`.

## Production on Cloud Run

Terraform lives in `infra/terraform/cloud-run` and provisions:

- an Artifact Registry Docker repository
- a Cloud Run service account
- a public Cloud Run service for the FlipCup image

The production workflow validates the app, bootstraps Artifact Registry via Terraform, builds the root image, pushes it, and applies the full Cloud Run configuration.

Required GitHub configuration:

- secret `GCP_PROJECT_ID`
- secret `GCP_WORKLOAD_IDENTITY_PROVIDER`
- secret `GCP_DEPLOYER_SERVICE_ACCOUNT`
- optional variable `GCP_REGION`
- optional variable `GCP_ARTIFACT_REPOSITORY`
- optional variable `GCP_CLOUD_RUN_SERVICE`

### Cloud Run websocket and state notes

Cloud Run supports WebSockets, but they are still long-lived HTTP requests. For this app that means:

- request timeout is set to `3600s`
- max instances stay at `1` so in-memory game state is not split across replicas
- websocket clients should be expected to reconnect after deploys or request timeout boundaries
- Cloud Run is not guaranteed to stay inside the free tier because active websocket connections keep the instance billable while they are open

Those defaults intentionally favor correctness over aggressive autoscaling.

## CI and deployment gates

Validation is centralized in `.github/workflows/validate.yml` and reused by:

- `CI` on every push and pull request
- `Deploy staging`
- `Deploy production`

That keeps Playwright, UI build, and Go test coverage consistent across normal development and deployments.

## Protecting `main`

Configure `main` in GitHub branch protection or rulesets with these minimum rules:

- require pull requests before merging
- block direct pushes
- require the `tests` status check from the `CI` workflow
- optionally require approvals before merge

The workflows in this repo assume deployment happens only after code reaches `main` through that PR flow.

## Legacy Fly.io support

`fly.toml` remains in the repository as an optional legacy target, but it now builds from the shared root `Dockerfile`. Cloud Run is the primary production path.
