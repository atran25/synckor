# Build stage
FROM golang:1.23-bookworm AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the go server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o synckor ./cmd/server/main.go

# Build the go migration binary
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go

# Download litestream
ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.8/litestream-v0.3.8-linux-amd64-static.tar.gz /tmp/litestream.tar.gz
RUN tar -C /usr/local/bin -xzf /tmp/litestream.tar.gz

# Final stage
FROM debian:bookworm
WORKDIR /app

# Install sqlite3
RUN apt-get update && apt-get install -y --no-install-recommends sqlite3

# Install ca-certificates and timezone data
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

# Copy the go binaries and the litestream binary
COPY --from=build /usr/local/bin/litestream /usr/local/bin/litestream
COPY --from=build /app/synckor /app/synckor
COPY --from=build /app/migrate /app/migrate

# Copy the api.yaml and the openapi docs
COPY /api.yaml /app/api.yaml
COPY /doc/index.html /app/doc/index.html

# Copy the configuration file and the run script
COPY /etc/litestream.yml /etc/litestream.yml
COPY /scripts/run.sh /app/scripts/run.sh

# Make sure the binary is executable
RUN chmod +x /app/synckor
RUN chmod +x /app/scripts/run.sh

RUN mkdir -p /app/db/sqlite
COPY /db/migrations /app/db/migrations
#COPY /db/sqlite /app/db/sqlite
ENTRYPOINT ["/app/scripts/run.sh"]