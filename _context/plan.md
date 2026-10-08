# Plan: Multi-Environment Deployment Strategy

## Problem
The current project has local Docker Compose development and a Fly.io production deployment path, but the desired target setup is:

- local development via Docker Compose
- staging on a local Raspberry Pi cluster running Nomad, Consul, and Traefik
- production on Google Cloud Run instead of Fly.io

The deployment story should make it easy to switch targets without maintaining three unrelated app packaging flows. In addition, Playwright tests should run on pushes to any branch and act as pre-deployment checks, and `main` should be protected so changes only land through pull requests from feature branches.

## Proposed Approach
Use a single app image and environment-specific deployment definitions:

- **Dev**: keep Docker Compose as the local developer workflow
- **Stage**: deploy the same app image to Nomad, fronted by Traefik, with environment values tuned for the local cluster
- **Prod**: deploy the same app image to Google Cloud Run using Terraform + GitHub Actions

The app packaging should be normalized so Fly-specific assumptions are removed or isolated. Environment selection should happen through deployment config and environment variables, not by rebuilding the app in completely different ways for each target. CI should validate branches continuously, and deployment workflows should depend on successful Playwright runs before stage or prod rollout.

## Todos
1. **Define environment contract**
   - Identify the runtime configuration the app needs per environment
   - Normalize variables such as public websocket/base URL configuration
   - Decide how a single image is configured differently for dev, stage, and prod

2. **Unify container build**
   - Replace or refactor the Fly-specific image path into a reusable deployment image
   - Preserve local Docker Compose development ergonomics
   - Verify the UI build/runtime configuration matches the intended environment

3. **Add production target for Cloud Run**
   - Add Terraform for Artifact Registry and Cloud Run
   - Configure Cloud Run for the app's in-memory / websocket constraints
   - Add GitHub Actions deployment flow for production
   - Make Fly.io optional or removable rather than the default prod target

4. **Add staging target for Nomad / Consul / Traefik**
   - Add Nomad job definitions for the app image
   - Define service registration and routing assumptions for Consul + Traefik
   - Document how staging differs from production

5. **Add CI / deployment gates**
   - Run Playwright tests on pushes to any branch
   - Reuse or require the Playwright workflow as a pre-deployment check for staging and production
   - Ensure failed Playwright runs block deployment jobs

6. **Protect the main branch**
   - Require pull requests for all changes into `main`
   - Prevent direct pushes and direct merges outside the PR workflow
   - Require the expected status checks before merge

7. **Document environment switching**
   - Explain local dev, stage deploy, and prod deploy flows
   - Document how to switch targets cleanly
   - Remove outdated Fly-first language once Cloud Run is the primary production path

## Notes / Considerations
- Local development already exists and should remain the easiest path for iteration.
- Staging and production should share the same application artifact as much as possible.
- The current UI env files appear inverted relative to their names:
  - `ui/.env` currently points at `flipcup.fly.dev`
  - `ui/.env.prod` currently points at `localhost:8080`
  This should be corrected during implementation as part of the environment contract work.
- Because game state is in memory, both Cloud Run and Nomad rollout/scaling settings should stay conservative unless shared state is introduced later.
- Branch protection may be enforced through GitHub branch protection rules or rulesets, but the plan should assume `main` becomes PR-only with required checks.
