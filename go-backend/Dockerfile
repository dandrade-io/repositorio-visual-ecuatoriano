# Stage 1: Build the Golang Application
ARG GOLANG_VERSION=1.27.1-alpine3.24

FROM golang:${GOLANG_VERSION} AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# Stage 2: Prepare server (distroless: CA certs + non-root user)
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /app/server /server

EXPOSE 8080
USER nonroot
ENTRYPOINT ["/server"]
