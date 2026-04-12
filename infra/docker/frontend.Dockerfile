FROM node:20-alpine AS dev

WORKDIR /app
COPY frontend/ ./

EXPOSE 3000
CMD ["sh", "-c", "if [ -f package.json ]; then npm install && npm run dev -- --host 0.0.0.0 --port 3000; else echo 'frontend/package.json not found; skipping frontend startup' && sleep infinity; fi"]

FROM node:20-alpine AS prod-builder

WORKDIR /app
COPY frontend/ ./
RUN if [ -f package.json ]; then npm install && npm run build; else mkdir -p .output/server .output/public && printf '%s\n' '<!doctype html><title>Frontend not configured</title><h1>Frontend not configured</h1>' > .output/public/index.html && printf '%s\n' 'console.log(\"Frontend not configured\")' > .output/server/index.mjs; fi

FROM node:20-alpine AS prod

WORKDIR /app
ENV HOST=0.0.0.0
ENV PORT=3000
COPY --from=prod-builder /app/.output ./.output
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
