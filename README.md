# TorrServer-Turbo (`torrserver-gst`)

Custom high-performance BitTorrent streaming server distribution for CineClaw with GStreamer 1.24 remuxing and instant time-to-first-frame (TTFF) optimizations.

## Features
- **TorrServer MatriX**: Enhanced BitTorrent streaming engine built on anacrolix/torrent.
- **GStreamer 1.24 Remuxing**: Integrated GStreamer pipeline for instant audio/subtitle track probing (`/gst/:hash/probe`) and HLS remuxing (`/torr/gst/<hash>/master.m3u8`).
- **Turbo BitTorrent Optimizations**:
  - Pre-resolved Tier-1 DHT bootstrap nodes in memory.
  - Aggressive peer discovery (40/40 dial limiters, 60 half-open conns, 200/50 peer watermarks).
  - 16-chunk initial pipeline per peer (pulling 6.4MB on initial round-trip).
  - 250ms chunk hedging watchdog duplicate-requesting stalled chunks (>p95 or >4s).
  - In-cache piece buffer free-list bypassing Go GC STW pauses.
  - Multi-season TV pack boundary piece range mapping.
- **Multi-Arch Docker Images**: Published to GHCR (`ghcr.io/cineclaw/torrserver-gst:latest`, `linux/amd64`, `linux/arm64`).

## Build
```bash
docker build -t ghcr.io/cineclaw/torrserver-gst:latest .
```
