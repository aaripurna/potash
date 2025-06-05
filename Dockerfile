FROM golang:1.24.3-alpine as go-builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o application


FROM oven/bun:latest as js-builder
WORKDIR /app
COPY package.json /app
COPY bun.lock /app
COPY vite.config.js /app
COPY assets /app/assets

RUN bun install
RUN bunx vite build


# Runtime
FROM scratch
WORKDIR /app
COPY --from=go-builder /app/application /app
COPY views /app/views
COPY public /app/public
COPY --from=js-builder /app/public/vite /app/public/vite

EXPOSE 80

ENTRYPOINT [ "/app/application" ]

