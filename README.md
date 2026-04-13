# DevHub

**DevHub** is an Internal Developer Platform (IDP) that streamlines the developer experience inside an organization. It allows engineers to **scaffold**, **deploy**, and **manage** cloud-native services with a unified interface.

> Built with Go, Vue.js, TypeScript, Kubernetes, Docker, and Python.


## 🔧 Features

- 🧱 Service scaffolding from templates (REST, gRPC, cron jobs, etc.)
- 🚀 One-click deployment to Kubernetes clusters
- 📦 CI/CD pipeline integration
- 📊 Dashboard for runtime metrics and deployment history
- 🔐 Role-based access controls for different teams
- 📁 Git integration (GitHub/GitLab/Bitbucket)
- 🧪 API testing and endpoint validation
- 🔄 Background job management
- 🧭 Plugin system for extending features

---

## 🛠 Tech Stack

| Layer         | Tech                                                                 |
|---------------|----------------------------------------------------------------------|
| Frontend      | Vue 3 + TypeScript + Headless UI + Tailwind CSS                      |
| Backend       | Go (REST API & scaffolding), Python (automation tasks & plugins)     |
| DevOps        | Docker, Kubernetes, Helm, GitHub Actions                             |
| CI/CD         | ArgoCD or GitHub Actions                                             |
| Storage       | PostgreSQL / Redis / Object Storage                                  |


## 🚀 Getting Started

```bash
# Clone the repo
git clone https://github.com/yourusername/devhub.git && cd devhub

# Prepare local environment
./scripts/bootstrap.sh

# Start the full local stack
./scripts/dev.sh up --build

# Run database migrations
./scripts/migrate.sh up

# Postgres is exposed on localhost:5433 by default to avoid
# colliding with an existing local Postgres on 5432
```


## 🧩 Project Structure

```markdown
devhub/
├── backend/                # Go backend service
│   ├── cmd/                # Application entry points (e.g. main.go)
│   ├── internal/           # Private application logic
│   │   ├── api/            # Route handlers (REST/gRPC)
│   │   ├── config/         # Configuration management
│   │   ├── domain/         # Core domain models and services
│   │   ├── infra/          # Infrastructure integrations (DB, external APIs)
│   │   ├── server/         # Server setup and lifecycle
│   │   ├── usecase/        # Business use cases
│   │   └── util/           # Utility functions
│   ├── migrations/         # Database schema migrations
│   ├── pkg/                # Shared public packages
│   ├── go.mod              # Go module definition
│   ├── go.sum              # Go module checksums
│   └── main.go             # Main application entry
│
├── frontend/               # Vue 3 + Tailwind dashboard
│   ├── src/
│   │   ├── components/     # Shared UI components
│   │   ├── layouts/        # Layout wrappers
│   │   ├── pages/          # Route-based views
│   │   ├── composables/    # Reusable logic (e.g. useFetch)
│   │   ├── stores/         # Pinia/Vuex state management
│   │   ├── assets/         # Static files, images
│   │   └── main.ts         # Entry point
│   └── vite.config.ts      # Build config
│
├── plugins/                # Optional Python plugin system
│   ├── scaffolders/        # Python service scaffolding logic
│   ├── runners/            # Background job executors
│   └── utils/              # Helpers for Python tasks
│
├── infra/                  # Infrastructure
│   ├── kubernetes/         # K8s manifests, Helm charts
│   ├── docker/             # Dockerfiles, entrypoints
│   └── terraform/          # Optional Terraform infra
│
├── scripts/                # Thin wrappers around common docker/go workflows
│   ├── dev.sh              # Start the local compose stack
│   ├── bootstrap.sh        # Create local env files for first run
│   ├── migrate.sh          # Run database migrations inside the backend container
│   ├── generate.sh         # Run backend code generation commands
│   └── docker-build-and-run.sh # Shared docker compose wrapper
│
├── templates/              # Service templates
│   ├── go-http/
│   ├── node-api/
│   └── python-worker/
│
├── workflows/                # CI/CD and automation scripts
│   ├── deploy.yaml           # Deployment workflow
│   ├── resource-provision.yaml # Infrastructure provisioning workflow
│   ├── rollback.yaml         # Rollback workflow
│   └── service-create.sh     # Service creation helper script
│
├── docs/                     # Markdown docs (API, onboarding, etc)
│   ├── architecture.md
│   ├── getting-started.md
│   └── roadmap.md
│
├── .github/                # GitHub Actions CI/CD workflows
│   └── workflows/
        ├── portal-ci.yaml
        ├── control-plane-ci.yaml
        ├── actions-ci.yaml
        └── infra-ci.yaml
│
├── docker-compose.yml      # Shared compose services and defaults
├── docker-compose.dev.yml  # Development overrides
├── docker-compose.prod.yml # Production overrides
├── README.md
└── LICENSE
```

## Helm And Argo CD

The repo now includes a deployable Helm chart at [infra/kubernetes/helm/devhub/values.yaml](/Users/hungpham/workspace/personal/devhub/infra/kubernetes/helm/devhub/values.yaml) and an Argo CD `Application` at [infra/kubernetes/argocd/devhub.yaml](/Users/hungpham/workspace/personal/devhub/infra/kubernetes/argocd/devhub.yaml).

The chart deploys:

- `api`: the Go HTTP server
- `worker`: a separate `Deployment` that runs `/app/devhub sync-worker`

For the fastest local end-to-end Argo CD smoke test, the included `Application` currently points to the public `argoproj/argocd-example-apps` repository over HTTPS and syncs the `guestbook` example. This avoids SSH repo credential setup while you verify worker-triggered syncs.

The backend worker now includes a real deployment runner that executes an Argo CD sync for each pending deployment. It invokes:

Example:

```bash
argocd app sync <service> --revision <version> --server <argocd-server> --auth-token <token>
```

It uses a hardcoded timeout of `10m`. The worker expects `ARGOCD_SERVER` and `ARGOCD_AUTH_TOKEN` in its runtime environment. The Helm chart sets `ARGOCD_SERVER` for the worker by default and exposes `secrets.argocdAuthToken`.

To run only the deployment worker:

```bash
go run ./backend sync-worker --types deployment
```

For local Docker Compose runs, the worker reads Argo CD credentials from `.env` via the `backend` service. Set:

```bash
ARGOCD_AUTH_TOKEN=<token>
```

After verifying the flow, update `spec.source.repoURL` and `spec.source.path` in [infra/kubernetes/argocd/devhub.yaml](/Users/hungpham/workspace/personal/devhub/infra/kubernetes/argocd/devhub.yaml) back to your real GitOps source.

Example:

```bash
helm upgrade --install devhub infra/kubernetes/helm/devhub \
  --namespace devhub \
  --create-namespace \
  --set image.repository=ghcr.io/phamphihungbk/devhub-backend \
  --set image.tag=latest \
  --set secrets.tokenSecret="$TOKEN_SECRET"

kubectl apply -f infra/kubernetes/argocd/devhub.yaml
```

### Local Argo CD UI With Minikube

If you want to use the Argo CD web UI locally, the repo now includes [scripts/argocd.sh](/Users/hungpham/workspace/personal/devhub/scripts/argocd.sh).

Typical flow:

```bash
minikube start
./scripts/argocd.sh all
```

That command will:

- install Argo CD into the current cluster
- apply the DevHub Argo CD `Application`
- start a local port-forward for the UI on `http://127.0.0.1:8081`
- print the default `admin` password from `argocd-initial-admin-secret`

The install step uses server-side apply to avoid the Kubernetes annotation-size error that can happen with Argo CD CRDs such as `applicationsets.argoproj.io`.

You can also run the steps individually:

```bash
./scripts/argocd.sh install
./scripts/argocd.sh app
./scripts/argocd.sh ingress
./scripts/argocd.sh domain
./scripts/argocd.sh configure
./scripts/argocd.sh ui
./scripts/argocd.sh password
./scripts/argocd.sh token
```

Or via Make:

```bash
make argocd-ui
make argocd-token
```

Note: the included Argo CD `Application` manifest points at `git@personal:phamphihungbk/devhub.git`. Make sure your Argo CD instance can reach that Git remote and has credentials configured for it.

## Local Gitea For Multi-Repo Testing

If you want a lightweight local Git server on your MacBook for testing many repositories, tags, and branches, Gitea is now included in the main Docker Compose stack and built from [infra/docker/gitea.Dockerfile](/Users/hungpham/workspace/personal/devhub/infra/docker/gitea.Dockerfile).

Start Gitea:

```bash
./scripts/dev.sh up -d gitea
```

Open the UI at [http://localhost:3000](http://localhost:3000).

To expose Gitea on the same local-domain setup as DevHub, run:

```bash
make setup-local-https
```

Then open [https://gitea.devhub.local](https://gitea.devhub.local).

Recommended first-run settings:

- database: `SQLite3`
- instance URL: `https://gitea.devhub.local/`
- SSH server domain: `gitea.devhub.local`
- SSH server port: `2222`

The compose file exposes:

- HTTP on `localhost:3000`
- SSH Git access on `localhost:2222`
- HTTPS via NGINX on `https://gitea.devhub.local`

Common commands:

```bash
./scripts/dev.sh ps gitea
./scripts/dev.sh logs -f gitea
./scripts/dev.sh stop gitea
```

To clone over HTTP after creating a repo:

```bash
git clone https://gitea.devhub.local/<user>/<repo>.git
```

To clone over SSH, add your SSH key in Gitea and use port `2222`:

```bash
git clone ssh://git@gitea.devhub.local:2222/<user>/<repo>.git
```

For Argo CD inside Minikube, prefer the HTTP repo URL using `host.minikube.internal`, for example:

```text
http://host.minikube.internal:3000/<user>/<repo>.git
```

To make scaffold requests publish directly into Gitea instead of returning a local `file://` path, set these values in [.env.example](/Users/hungpham/workspace/personal/devhub/.env.example) / `.env`:

```bash
GITEA_URL=http://gitea:3000
GITEA_EXTERNAL_URL=https://gitea.devhub.local
GITEA_USERNAME=<your-gitea-user>
GITEA_TOKEN=<your-gitea-token>
GITEA_OWNER=<optional-org-or-user>
```

When `GITEA_USERNAME` and `GITEA_TOKEN` are present, scaffold plugins will:

- generate the service locally
- create a repository in Gitea
- commit the generated files
- push the initial `main` branch
- store the Gitea clone URL as the scaffold result

To generate a token for the deployment worker after the Argo CD server is reachable locally:

```bash
./scripts/argocd.sh token
```

The helper ensures the `admin` account has `apiKey, login` enabled and grants local admin RBAC before generating the token. It prints an `export ARGOCD_AUTH_TOKEN=...` line you can paste into your shell or `.env`.

The local HTTPS helper [scripts/setup-local-https.sh](/Users/hungpham/workspace/personal/devhub/scripts/setup-local-https.sh) now also updates `/etc/hosts` for `argocd.devhub.local` when Minikube is running. You can override the detected IP with `DEVHUB_ARGOCD_IP`.

### Domain Access Through NGINX Ingress

For Minikube, the repo also includes [infra/kubernetes/argocd/argocd-ui-ingress.yaml](/Users/hungpham/workspace/personal/devhub/infra/kubernetes/argocd/argocd-ui-ingress.yaml), which exposes the Argo CD UI through the NGINX ingress addon at `argocd.devhub.local`.

Run:

```bash
./scripts/argocd.sh ingress
./scripts/argocd.sh domain
```

Then add the printed `minikube ip` entry to your local `/etc/hosts`, for example:

```text
192.168.49.2 argocd.devhub.local
```

After that, open:

```text
https://argocd.devhub.local
```

The ingress forwards traffic to the Argo CD server service over HTTPS, while keeping browser access simple for local development.

## 🛣️ `ROADMAP.md`

```markdown
# 📍 DevHub Roadmap

This roadmap outlines the key milestones for DevHub from MVP to full internal platform.

---

## ✅ Phase 1: MVP (Core Features)
- [x] Vue + Tailwind UI dashboard
- [x] Go-based backend API
- [x] Scaffold service templates (REST, cron jobs)
- [x] Kubernetes deployment integration
- [x] Dockerized frontend/backend
- [x] Local dev environment (Docker Compose)

---

## 🚧 Phase 2: Developer Experience
- [ ] Add form-based UI for service scaffolding
- [ ] GitHub/GitLab repo scaffolding & commit hooks
- [ ] Deployment logs + terminal access
- [ ] Add background job template
- [ ] JWT or OAuth2 authentication

---

## 🔜 Phase 3: DevOps Automation
- [ ] CI/CD pipeline templates (GitHub Actions, ArgoCD)
- [ ] Service status overview panel (health checks)
- [ ] Helm chart management UI
- [ ] Automatic rollback on failed deploy

---

## 🧠 Phase 4: Extensibility & Insights
- [ ] Plugin system (Python modules or Webhooks)
- [ ] Metrics via Prometheus + Grafana
- [ ] Usage tracking (most active projects, teams)
- [ ] Notifications (Slack, Email, Discord)

---

## 🧪 Future Ideas
- [ ] Internal ChatGPT plugin integration
- [ ] AI-assisted scaffold suggestion
- [ ] Secret manager integration (Vault / SOPS)
- [ ] Feature flag UI


---


┌───────────────────────────────────────────────────────────────────┐
│                           FRONTEND UI                             │
│     Create Service | Deploy | View Metrics | Manage Plugins       │
└───────────────────────────────┬───────────────────────────────────┘
                                │
                                ▼
┌───────────────────────────────────────────────────────────────────┐
│                      GO CONTROL PLANE API                         │
│                                                                   │
│  Scaffold API   Deploy API   Metrics API   RBAC   Plugin API      │
└───────────────┬──────────────┬─────────────┬──────────────────────┘
                │              │             │
                ▼              ▼             ▼
         ┌────────────┐  ┌────────────┐  ┌──────────────┐
         │ PostgreSQL │  │   Worker   │  │ Plugin Reg.  │
         │            │  │   System   │  │              │
         └─────┬──────┘  └─────┬──────┘  └──────┬───────┘
               │               │                │
               ▼               ▼                ▼
     scaffold_requests   ScaffoldWorker     Scaffold Plugins
     deployments         DeploymentWorker   Deploy Plugins
     test jobs           TestWorker         Test Plugins
     audit/history       PluginWorker       Integration Plugins
               │               │                │
               └───────────────┴────────────────┘
                               │
                               ▼
                    External Systems / Runtime
            Git + CI/CD + Kubernetes + Metrics + APIs
