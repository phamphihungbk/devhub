ARG MODE

FROM node:20-alpine AS base

WORKDIR /app

COPY frontend/ ./

# RUN npm install

# Production build stage
FROM base AS prod

# RUN npm run build

# Development build stage
FROM base AS dev

# Final image
FROM nginx:alpine AS final

# COPY --from=prod /app/dist /usr/share/nginx/html

EXPOSE 80

EXPOSE 3000

CMD ["sh", "-c", "if [ '$MODE' = 'dev' ]; then cd /app && npm run dev; else nginx -g 'daemon off;'; fi"]
