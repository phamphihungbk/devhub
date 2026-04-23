# API Docs

## Overview

The DevHub API is organized around the control-plane lifecycle:

- approval policies and approval requests
- projects
- services
- scaffold requests
- releases
- deployments
- plugins
- auth and users

The backend is implemented in Go under [`backend/internal/api/http`](../backend/internal/api/http).

## Main Resource Flow

```text
Project
  -> Service
    -> Release
      -> Deployment
```

Scaffold requests belong to a project and are one path to creating a new service record.

## Core Endpoints

### Auth

- `POST /auth/login`
- `POST /auth/logout`
- `GET /auth/me`

### Approval

- `POST /approval-policies`
- `GET /approval-requests`
- `POST /approval-requests/:approval-request/decisions`

Approval decisions accept both a decision and a comment:

```json
{
  "decision": "approve",
  "comment": "Validated rollout prerequisites and approved for execution."
}
```

### Projects

- `GET /projects`
- `POST /projects`
- `GET /projects/:project`
- `PATCH /projects/:project`
- `DELETE /projects/:project`

### Services

- `GET /projects/:project/services`

Service state is consumed by the frontend service details page and by service-scoped release and deployment flows.

### Teams

- `GET /teams`
- `POST /teams`
- `PATCH /teams/:team`

### Scaffold Requests

- `GET /projects/:project/scaffold-requests`
- `POST /projects/:project/scaffold-requests`
- `GET /scaffold-requests/:scaffoldRequest`
- `DELETE /scaffold-requests/:scaffoldRequest`

Important request shape:

```json
{
  "plugin_id": "plugin-id",
  "environment": "dev",
  "variables": {
    "service_name": "payment",
    "module_path": "github.com/acme/payment",
    "port": 8080,
    "database": "postgres",
    "enable_logging": true
  }
}
```

### Releases

- `GET /services/:service/releases`
- `POST /services/:service/releases`

Example create payload:

```json
{
  "plugin_id": "releaser-plugin-id",
  "tag": "v1.0.0",
  "target": "main",
  "name": "Payment v1.0.0",
  "notes": "Initial release"
}
```

### Deployments

- `GET /services/:service/deployments`
- `POST /services/:service/deployments`
- `GET /deployments/:deployment`
- `PATCH /deployments/:deployment`
- `DELETE /deployments/:deployment`

Example create payload:

```json
{
  "plugin_id": "deployer-plugin-id",
  "environment": "staging",
  "version": "v1.0.0"
}
```

Deployments are service-scoped and version-based. In the frontend, `version` is selected from the service’s releases.

### Plugins

- `GET /plugins`
- `POST /plugins`
- `GET /plugins/:plugin`
- `PATCH /plugins/:plugin`
- `DELETE /plugins/:plugin`

Plugin types in active use:

- `scaffolder`
- `releaser`
- `deployer`

## Frontend Mapping

The Vue UI uses these route-level flows:

- `Dashboard` -> team overview and quick health snapshot
- `Approvals` -> review approval requests and leave decision comments
- `Projects` -> list and filter projects
- `Services` -> browse services across projects
- `Releases` -> browse release history and release timeline
- `Project details` -> inspect services and scaffold
- `Service details` -> create release and deploy selected release versions
- `Plugins` -> inspect plugin registry

Frontend API clients live under [`frontend/src/services/api`](../frontend/src/services/api).

## Notes

- Release records are service-scoped.
- Deployments are service-scoped.
- Deployment `version` should match a release tag for that service.
- Worker payloads are built in Go and executed through Python plugins.
