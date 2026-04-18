FROM node:20-alpine AS dev

WORKDIR /app
COPY frontend/ ./
RUN corepack enable

EXPOSE 3000
CMD ["sh", "-c", "if [ -f package.json ]; then pnpm install && pnpm dev --host 0.0.0.0 --port 3000; else echo 'frontend/package.json not found; skipping frontend startup' && sleep infinity; fi"]

FROM node:20-alpine AS prod-builder

WORKDIR /app
COPY frontend/ ./
RUN corepack enable
RUN if [ -f package.json ]; then pnpm install && pnpm build; else mkdir -p dist && printf '%s\n' '<!doctype html><title>Frontend not configured</title><h1>Frontend not configured</h1>' > dist/index.html; fi

FROM node:20-alpine AS prod

WORKDIR /app
ENV HOST=0.0.0.0
ENV PORT=3000
RUN corepack enable
COPY frontend/package.json ./package.json
COPY --from=prod-builder /app/node_modules ./node_modules
COPY --from=prod-builder /app/dist ./dist
EXPOSE 3000
CMD ["sh", "-c", "pnpm preview --host 0.0.0.0 --port 3000"]
