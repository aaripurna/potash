FROM oven/bun:latest as js-builder
WORKDIR /app
COPY . .
RUN bun install
RUN NODE_ENV=production bunx vite build

FROM golang:1.25.1-alpine as go-builder
WORKDIR /app
COPY . .
RUN go mod download
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

