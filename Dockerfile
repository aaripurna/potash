FROM oven/bun:latest as js-install
RUN mkdir -p /temp/prod
COPY package.json bun.lock /temp/prod/
RUN cd /temp/prod && NODE_ENV=production bun install --frozen-lockfile --production

FROM oven/bun:latest as js-builder
WORKDIR /app
COPY . .
COPY --from=js-install /temp/prod/node_modules node_modules
RUN NODE_ENV=production bunx vite build

FROM golang:1.25.1-alpine as go-builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x
ENV GOCACHE=/root/.cache/go-build
COPY --from=js-builder /app/dist /app/public
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o application

# Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=go-builder /app/application /app
COPY --from=js-builder /app/dist /app/public
COPY views /app/views
EXPOSE 80

ENTRYPOINT [ "/app/application" ]

