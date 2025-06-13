# https://github.com/GoogleCloudPlatform/golang-samples/blob/main/run/helloworld/Dockerfile
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o main ./cmd

FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main /app/main
# Move the SQL migration files into the Dockerfile so that
# we can apply database migrations using the Docker container if needed.
COPY --from=builder /app/db/migrations /app/db/migrations

USER 1001

CMD ["/app/main", "server"]