FROM node:20-alpine AS dev

WORKDIR /app
COPY frontend/ ./

EXPOSE 3000
CMD ["sh", "-c", "if [ -f package.json ]; then npm install && npm run dev -- --host 0.0.0.0 --port 3000; else echo 'frontend/package.json not found; skipping frontend startup' && sleep infinity; fi"]

FROM node:20-alpine AS prod-builder

WORKDIR /app
COPY frontend/ ./
RUN if [ -f package.json ]; then npm install && npm run build; else mkdir -p dist && printf '%s\n' '<!doctype html><title>Frontend not configured</title><h1>Frontend not configured</h1>' > dist/index.html; fi

FROM nginx:1.27-alpine AS prod

COPY --from=prod-builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
