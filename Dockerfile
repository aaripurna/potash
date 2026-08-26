# syntax=docker/dockerfile:1

FROM oven/bun:latest AS js-builder
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY . .
RUN bunx vite build

FROM golang:1.27.0-alpine AS go-builder
ARG GOOS=linux
ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
# public/ and views/ are embedded into the binary, so dist must land before the
# build - nothing is copied into the runtime stage.
COPY --from=js-builder /app/dist /app/public
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS="${GOOS}" GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o application

# Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata \
    && adduser -D -u 10001 app
WORKDIR /app
ENV NODE_ENV=production
# /usr/local/bin is already on PATH, so the binary is invocable by name.
COPY --from=go-builder /app/application /usr/local/bin/application
USER app

ENTRYPOINT [ "application" ]
CMD [ "serve", "--bind", "0.0.0.0" ]
