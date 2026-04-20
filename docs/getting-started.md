# Getting Started

## Local Setup

```bash
git clone https://github.com/yourusername/devhub.git
cd devhub
./scripts/bootstrap.sh
./scripts/dev.sh up --build
./scripts/migrate.sh up
```

Notes:

- Postgres is exposed on `localhost:5433` by default.
- The generated `.env` is the main place to configure backend, worker, SCM, and Argo CD integration.

Helpful local domains:

- DevHub UI: [https://devhub.local](https://devhub.local)
- DevHub API: [https://api.devhub.local](https://api.devhub.local)
- Gitea UI: [https://gitea.devhub.local](https://gitea.devhub.local)
- Argo CD UI: [https://argocd.devhub.local](https://argocd.devhub.local)

To configure local host mappings and trust the generated certificate on macOS:

```bash
make setup-local-https
```

## Watch Mode

Run the backend locally with file watch while keeping support services in Docker:

```bash
make backend-watch
```

Run the frontend in watch mode:

```bash
make frontend-watch
```

## Minikube With Local Registry

To connect the local `devhub-registry` to Minikube, use:

```bash
make minikube-registry
```

That target will:

- start or recreate `devhub-registry`
- recreate the Minikube cluster
- start Minikube with `host.minikube.internal:5001` configured as an insecure registry
- verify that Minikube can reach the registry

After that:

- push from the host to `host.docker.internal:5001/<image>:<tag>`
- use `host.minikube.internal:5001/<image>:<tag>` inside Kubernetes or GitOps values

Example:

```bash
docker tag supperfast:latest host.docker.internal:5001/supperfast:latest
docker push host.docker.internal:5001/supperfast:latest
```

```yaml
image:
  repository: host.minikube.internal:5001/supperfast
  tag: latest
```

## Argo CD And GitOps

DevHub includes:

- a platform Helm chart at [`infra/kubernetes/helm/devhub/values.yaml`](../infra/kubernetes/helm/devhub/values.yaml)
- an Argo CD `ApplicationSet` at [`infra/kubernetes/argocd/devhub.yaml`](../infra/kubernetes/argocd/devhub.yaml)

The worker executes deployment jobs by syncing Argo CD applications with the selected release version.

Example sync shape:

```bash
argocd app sync <service> --revision <version> --server <argocd-server> --auth-token <token>
```

Required worker/runtime env:

```bash
ARGOCD_SERVER=<argocd-server>
ARGOCD_AUTH_TOKEN=<token>
```

The generated GitOps values files are expected to look like:

```yaml
appProject: "default"
appName: "payment"
appEnvironment: "dev"
appRepoURL: "http://host.minikube.internal:3000/phamphihungbk/payment.git"
appTargetRevision: "main"
appNamespace: "devhub"

nameOverride: "payment"
fullnameOverride: "payment"
```

The `ApplicationSet` watches environment files and creates one Argo CD application per file.

## Local Argo CD UI

Typical flow:

```bash
minikube start
./scripts/argocd.sh install
./scripts/argocd.sh repo
./scripts/argocd.sh app
./scripts/argocd.sh ui
```

Useful commands:

```bash
./scripts/argocd.sh password
./scripts/argocd.sh token
./scripts/argocd.sh ingress
./scripts/argocd.sh domain
```

Make aliases:

```bash
make argocd-ui
make argocd-token
```

## Local Gitea

Start Gitea:

```bash
./scripts/dev.sh up -d gitea
```

Open:

- [http://localhost:3000](http://localhost:3000)
- or [https://gitea.devhub.local](https://gitea.devhub.local) after:

```bash
make setup-local-https
```

Recommended first-run settings:

- database: `SQLite3`
- instance URL: `https://gitea.devhub.local/`
- SSH server domain: `gitea.devhub.local`
- SSH server port: `2222`

Clone examples:

```bash
git clone https://gitea.devhub.local/<user>/<repo>.git
git clone ssh://git@gitea.devhub.local:2222/<user>/<repo>.git
```

For Argo CD inside Minikube, prefer an internal HTTP URL such as:

```text
http://host.minikube.internal:3000/<user>/<repo>.git
```

## SCM Publishing For Scaffold Jobs

To publish scaffolded services directly into Gitea, set these values in `.env`:

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
- store the clone URL as the scaffold result
