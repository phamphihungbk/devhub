# Plugin Scaffolders

This folder contains local plugin entrypoints that the DevHub worker can execute.

Available scaffold_request:

- `plugins/scaffold_request/go_http_api/run.py`: Go HTTP API
- `plugins/scaffold_request/go_grpc_service/run.py`: Go gRPC service
- `plugins/scaffold_request/go_worker/run.py`: Go worker
- `plugins/scaffold_request/python_fastapi/run.py`: Python FastAPI service
- `plugins/scaffold_request/python_worker/run.py`: Python worker
- `plugins/scaffold_request/node_express_api/run.py`: Node.js Express API
- `plugins/scaffold_request/node_worker/run.py`: Node.js worker
- `plugins/scaffold_request/react_vite_app/run.py`: React Vite frontend
- `plugins/scaffold_request/vue_vite_app/run.py`: Vue Vite frontend
- `plugins/scaffold_request/nextjs_app/run.py`: Next.js frontend

Each scaffolder reads a JSON payload from stdin and prints a JSON result to stdout.

To scaffold a new local plugin folder:

```bash
./scripts/create-plugin.sh --name my-plugin --type scaffolder --description "My plugin"
```
