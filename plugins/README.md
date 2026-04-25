# Plugin Scaffolders

This folder contains local plugin entrypoints that the DevHub worker can execute.

Available scaffolders:

- `plugins/scaffolders/go_http_api/run.py`: Go HTTP API
- `plugins/scaffolders/go_grpc_service/run.py`: Go gRPC service
- `plugins/scaffolders/go_worker/run.py`: Go worker
- `plugins/scaffolders/python_fastapi/run.py`: Python FastAPI service
- `plugins/scaffolders/python_worker/run.py`: Python worker
- `plugins/scaffolders/node_express_api/run.py`: Node.js Express API
- `plugins/scaffolders/node_worker/run.py`: Node.js worker
- `plugins/scaffolders/react_vite_app/run.py`: React Vite frontend
- `plugins/scaffolders/vue_vite_app/run.py`: Vue Vite frontend
- `plugins/scaffolders/nextjs_app/run.py`: Next.js frontend

Each scaffolder reads a JSON payload from stdin and prints a JSON result to stdout.

To scaffold a new local plugin folder:

```bash
./scripts/create-plugin.sh --name my-plugin --type scaffolder --description "My plugin"
```
