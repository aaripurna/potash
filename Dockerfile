# syntax=docker/dockerfile:1

FROM oven/bun:latest AS js-builder
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY . .
RUN bunx vite build

FROM golang:1.27.0-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
# dist is embedded (//go:embed public/*) so it must land before the build
COPY --from=js-builder /app/dist /app/public
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o application

# Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata \
    && adduser -D -u 10001 app
WORKDIR /app
ENV NODE_ENV=production
COPY --from=go-builder /app/application /app
COPY --from=js-builder /app/dist /app/public
COPY views /app/views
USER app

ENTRYPOINT [ "/app/application" ]
CMD [ "serve", "--bind", "0.0.0.0" ]
