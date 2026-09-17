# APT Proxy

[![Docker Pulls](https://img.shields.io/docker/pulls/lj020326/apt-proxy.svg)](https://hub.docker.com/r/lj020326/apt-proxy)
[![Security Scan](https://github.com/lj020326/apt-proxy/actions/workflows/scan.yml/badge.svg)](https://github.com/lj020326/apt-proxy/actions/workflows/scan.yml)
[![Release](https://github.com/lj020326/apt-proxy/actions/workflows/release.yml/badge.svg)](https://github.com/lj020326/apt-proxy/actions/workflows/release.yml)

A lightweight (**< 10MB**), high-performance caching proxy for Linux package managers (APT, YUM/DNF, APK). It accelerates package installation by caching packages locally, drastically reducing build times and bandwidth consumption across server fleets and container build pipelines.

---

## Features

* **Multi-Distribution Support**: Works out of the box with Ubuntu/Debian (`apt`), CentOS (`yum`/`dnf`), and Alpine (`apk`).
* **Minimal Footprint**: Binary size under 10MB built with Go and Fiber.
* **Smart Mirror Selection**: Automatically benchmarks and routes requests to the fastest upstream mirror on startup.
* **Flexible Storage**: Supports local disk caching or offloading cache objects to S3-compatible storage (AWS S3, MinIO, OtterIO, Ceph, Cloudflare R2).
* **Observability**: Built-in `/healthz`, `/livez`, `/readyz`, and Prometheus `/metrics` endpoints.
* **Management API**: Protected REST endpoints (`/api/cache/stats`, `/api/cache/purge`, `/api/cache/cleanup`, `/api/mirrors/refresh`).

---

## Supported Architectures

Multi-arch Docker images are pushed to Docker Hub (`lj020326/apt-proxy`) and GitHub Container Registry (`ghcr.io/lj020326/apt-proxy`):

* `linux/amd64`
* `linux/arm64`
* `linux/arm/v7`

---

## Quick Start

### 1. Run Container

Run APT Proxy using a persistent volume for the local cache directory:

```bash
docker run -d \
  --name apt-proxy \
  -p 3142:3142 \
  -v apt-proxy-cache:/var/cache/apt-proxy \
  lj020326/apt-proxy:latest
```

### 2. Configure Clients

#### Ubuntu / Debian
Pass the proxy via environment variable during package operations:

```bash
export http_proxy=http://<HOST_IP>:3142
apt-get update && apt-get install -y vim
```

#### Dockerfile Usage
Speed up package installs inside your image builds:

```dockerfile
FROM ubuntu:24.04
ARG APT_PROXY
ENV http_proxy=${APT_PROXY}
RUN apt-get update && apt-get install -y \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*
```

Build passing your proxy container's address:
```bash
docker build --build-arg APT_PROXY=[http://172.17.0.1:3142](http://172.17.0.1:3142) -t my-app .
```

---

## Configuration

Configuration can be passed via CLI arguments, environment variables, or a mounted YAML configuration file (`/etc/apt-proxy/apt-proxy.yaml`).

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `APT_PROXY_HOST` | `0.0.0.0` | Network interface to bind |
| `APT_PROXY_PORT` | `3142` | Server listening port |
| `APT_PROXY_MODE` | `all` | Distribution mode (`all`, `ubuntu`, `debian`, `centos`, `alpine`) |
| `APT_PROXY_CACHEDIR` | `/var/cache/apt-proxy` | Local cache directory |
| `APT_PROXY_CACHE_MAX_SIZE` | `10` | Max cache size in GB (`0` disables LRU eviction) |
| `APT_PROXY_CACHE_TTL` | `168` | Cache TTL in hours (`0` disables TTL eviction) |
| `APT_PROXY_API_KEY` | *None* | API key for protected endpoints (enables API auth) |
| `APT_PROXY_STORAGE_BACKEND` | `disk` | Storage backend (`disk` or `s3`) |

### S3 Backend Configuration

To share cache across container replicas or outlive ephemeral nodes, point storage to an S3-compatible backend:

```bash
docker run -d \
  --name apt-proxy \
  -p 3142:3142 \
  -e APT_PROXY_STORAGE_BACKEND=s3 \
  -e APT_PROXY_S3_ENDPOINT=minio.example.com:9000 \
  -e APT_PROXY_S3_BUCKET=apt-cache \
  -e APT_PROXY_S3_ACCESS_KEY=your-access-key \
  -e APT_PROXY_S3_SECRET_KEY=your-secret-key \
  -e APT_PROXY_S3_USE_PATH_STYLE=true \
  lj020326/apt-proxy:latest
```

---

## Docker Compose Example

```yaml
version: '3.8'

services:
  apt-proxy:
    image: lj020326/apt-proxy:latest
    container_name: apt-proxy
    restart: unless-stopped
    ports:
      - "3142:3142"
    environment:
      - TZ=UTC
      - APT_PROXY_CACHE_MAX_SIZE=20
      - APT_PROXY_CACHE_TTL=168
    volumes:
      - apt-proxy-cache:/var/cache/apt-proxy
      - /etc/localtime:/etc/localtime:ro

volumes:
  apt-proxy-cache:
```

---

## Links & Resources

* **GitHub Repository**: [lj020326/apt-proxy](https://github.com/lj020326/apt-proxy)
* **Issue Tracker**: [GitHub Issues](https://github.com/lj020326/apt-proxy/issues)
* **Full Documentation**: [README.md](https://github.com/lj020326/apt-proxy#readme)
* **License**: Apache-2.0
