FROM nginx:1.27-alpine

COPY infra/nginx/default.conf.template /etc/nginx/templates/default.conf.template
