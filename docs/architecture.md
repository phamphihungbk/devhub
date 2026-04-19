# Architecture

## Overview

DevHub is split into four primary layers:

1. `frontend`
   - Vue 3 control plane UI
   - project, service, release, and deployment workflows
2. `backend`
   - Go HTTP API
   - usecase and domain layers
   - worker runners for async jobs
3. `plugins`
   - Python scaffolders, releasers, and deployers
   - executed by backend workers
4. `infra`
   - Helm chart for the platform
   - Argo CD `ApplicationSet`
   - Docker and local environment assets

## High-Level Flow

```text
Frontend UI
  -> Go HTTP API
  -> PostgreSQL / Redis
  -> Background workers
  -> Python plugins
  -> SCM / GitOps repo / Argo CD / Kubernetes
```

## Frontend

The admin console is route-first and centered around these flows:

- `Projects`
  - list projects
  - filter by environment, status, owner team
- `Project details`
  - inspect services
  - open scaffold request modal
  - review recent releases and deployments
- `Service details`
  - create release
  - select a release
  - deploy based on the selected release version
  - inspect deployments filtered by release tag
- `Plugins`
  - browse plugin registry by type, runtime, and scope

## Backend

The Go backend follows a layered structure:

- `internal/api/http`
  - HTTP handlers and routes
- `internal/usecase`
  - business workflows
- `internal/domain`
  - entities and repository interfaces
- `internal/infra`
  - database repositories
  - workers
  - external integrations

Key domains currently include:

- auth
- users
- projects
- services
- scaffold requests
- releases
- deployments
- plugins

## Workers And Plugins

Async job types map to plugin families:

- scaffold requests -> scaffold plugins
- releases -> releaser plugins
- deployments -> deployer plugins

The worker builds a payload, executes the plugin, then persists the result back into DevHub state.

## GitOps And Delivery

Scaffold jobs can:

- create the service repository
- generate service files
- push the initial branch
- create or update the GitOps values file

Release jobs work at the service level and typically produce a tag or SCM release artifact.

Deployment jobs are now service-scoped and deploy a selected release version. The deployment worker syncs Argo CD using the chosen version/tag.

## Repo Structure

```text
devhub/
├── backend/
├── frontend/
├── plugins/
├── infra/
├── scripts/
└── docs/
```
