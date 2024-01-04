FROM --platform=linux/amd64 caddy:2.4.6-alpine

COPY Caddyfile.production /etc/caddy/Caddyfile
