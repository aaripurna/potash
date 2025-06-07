FROM oven/bun:latest as js-builder
WORKDIR /app
COPY package.json /app
COPY bun.lock /app
COPY vite.config.js /app
COPY public /app/public
COPY assets /app/assets
RUN bun install
RUN bunx vite build

FROM golang:1.24.3-alpine as go-builder
WORKDIR /app
COPY . .
RUN go mod download
COPY --from=js-builder /app/dist /app/public
RUN CGO_ENABLED=0 GOOS=linux go build -o application

# Runtime
FROM scratch
WORKDIR /app
COPY --from=go-builder /app/application /app
COPY --from=js-builder /app/dist /app/public
COPY views /app/views
EXPOSE 80

ENTRYPOINT [ "/app/application" ]

