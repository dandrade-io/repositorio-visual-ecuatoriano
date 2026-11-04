ARG NODE_VERSION=24.21-alpine
ARG NGINX_VERSION=alpine3.22

## STAEGE 1: Build the Angular Application
FROM node:${NODE_VERSION} AS builder

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./

RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile

COPY . .

RUN pnpm run build


## STAGE 2: Prepare Nginx to Serve Static Files
FROM nginxinc/nginx-unprivileged:${NGINX_VERSION} AS runner

COPY nginx.conf /etc/nginx/nginx.conf
COPY --chown=nginx:nginx --from=builder /app/dist/*/browser /usr/share/nginx/html
COPY --chown=nginx:nginx public/assets/env.template.json /usr/share/nginx/html/assets/env.template.json
COPY --chmod=755 docker/entrypoint.sh /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
