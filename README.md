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

# Install frontend
cd frontend && npm install && npm run dev

# Install backend (Go)
cd ../backend && go run main.go

# Or run with Docker
docker-compose up --build
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
├── portal/                 # Vue 3 + Tailwind dashboard
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
├── scripts/                # Dev and setup scripts
│   ├── dev.sh              # Start local development environment
│   ├── bootstrap.sh        # Initial project setup and dependencies
│   ├── migrate.sh          # Run database migrations
│   └── generate.sh         # Code scaffolding helper
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
├── docker-compose.yml      # Fullstack local setup
├── README.md
└── LICENSE
```

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
