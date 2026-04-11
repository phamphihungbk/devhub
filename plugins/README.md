# Plugin Scaffolders

This folder contains local plugin entrypoints that the DevHub worker can execute.

Available scaffolders:

- `plugins/scaffolders/go_http_api/action.py`: minimal Go HTTP API
- `plugins/scaffolders/go_service/action.py`: minimal Go service/worker
- `plugins/scaffolders/node_http_api/action.py`: minimal Node.js HTTP API
- `plugins/scaffolders/python_worker/action.py`: minimal Python worker

Each scaffolder reads a JSON payload from stdin and prints a JSON result to stdout.

To scaffold a new local plugin folder:

```bash
./scripts/create-plugin.sh --name my-plugin --type scaffolder --description "My plugin"
```
