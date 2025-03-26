FROM --platform=$BUILDPLATFORM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    GOOS=$TARGETOS GOARCH=$GOARCH make

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/bin/redis-vs-nats /app/exec

ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT [ "/app/exec" ]

# default to redis
CMD [ "redis" ]
