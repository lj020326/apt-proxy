# specify-mirrors

Same as [`../basic/`](../basic/), but pins upstream Ubuntu and Debian mirrors via CLI flags. Useful when geo-detection picks a slow mirror, or when you want a fixed, reproducible upstream.

## What's different

The compose file passes flags to the container:

```yaml
command: --ubuntu=us:archive --debian=us:mit
```

Both `us:archive` and `us:mit` (and a handful of other shortcuts) are built-in aliases. You can also pass full URLs:

```yaml
command: >-
  --ubuntu=[https://archive.ubuntu.com/ubuntu/](https://archive.ubuntu.com/ubuntu/)
  --debian=[https://mirrors.mit.edu/debian/](https://mirrors.mit.edu/debian/)
```

## Run

```bash
cd examples/specify-mirrors
docker compose up -d
```

## Verify

```bash
curl -I http://127.0.0.1:3142/healthz

curl -x http://127.0.0.1:3142 \
     http://archive.ubuntu.com/ubuntu/dists/noble/Release
```

## Configuring via YAML instead of flags

The same overrides can live in [`../config-template/apt-proxy.yaml`](../config-template/apt-proxy.yaml) under the `mirrors:` block — handy if you also want to tune cache size, TLS, or storage backend at the same time.
