# Stage 1: Build TorrServer-Turbo with GStreamer support
FROM golang:1.25-bookworm AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

# Copy source trees
COPY anacrolix-torrent ./anacrolix-torrent
COPY server ./server

WORKDIR /build/server

# Build static binary with gst tags
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} CGO_ENABLED=0 go build \
    -ldflags="-s -w -checklinkname=0" \
    -tags=nosqlite,gst \
    -trimpath \
    -o /usr/bin/torrserver ./cmd

# Stage 2: Runtime environment with Ubuntu 24.04 and GStreamer 1.24
FROM ubuntu:24.04

LABEL maintainer="CineClaw <https://github.com/cineclaw>"

ENV DEBIAN_FRONTEND=noninteractive
ENV TS_PORT=8090
ENV TS_PATH=/opt/ts/db
ENV TS_TORRENTSDIR=/opt/ts/torrents
ENV TS_DONTKILL=1
ENV GODEBUG=madvdontneed=1

# Install runtime dependencies for TorrServer and GStreamer 1.0 plugins
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    libgstreamer1.0-0 \
    gstreamer1.0-tools \
    gstreamer1.0-plugins-base-apps \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav \
    && rm -rf /var/lib/apt/lists/*

# Copy our compiled TorrServer-Turbo binary from the builder
COPY --from=builder /usr/bin/torrserver /usr/bin/torrserver
RUN chmod +x /usr/bin/torrserver

COPY start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 8090

ENTRYPOINT ["/start.sh"]
