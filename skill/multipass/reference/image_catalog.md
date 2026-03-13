# Multipass Image Catalog

## Discovering Available Images

```bash
multipass find                   # List all available images
multipass find --format json     # JSON output
```

With MCP: `multipass_find_images` returns the full catalog as JSON.

## Ubuntu LTS Releases

These are stable, long-term support releases. Prefer these for most use cases.

| Alias | Release | Codename | Notes |
|-------|---------|----------|-------|
| `24.04` | Ubuntu 24.04 LTS | Noble Numbat | **Current default** — latest LTS |
| `22.04` | Ubuntu 22.04 LTS | Jammy Jellyfish | Previous LTS, widely supported |
| `20.04` | Ubuntu 20.04 LTS | Focal Fossa | Older LTS, still maintained |

To launch a specific version:
```bash
multipass launch 24.04 --name my-vm
multipass launch 22.04 --name legacy-vm
```

## Daily Builds

Unstable, pre-release images. Use only for testing upcoming Ubuntu features.

```bash
multipass launch daily:24.10 --name daily-vm
```

## Core Images

Minimal Ubuntu Core images (snap-based). Smaller footprint, fewer packages.

```bash
multipass launch core22 --name core-vm
```

## Blueprints

Pre-configured images for specific workloads. Availability varies by platform.

```bash
multipass launch docker --name docker-vm     # Docker pre-installed (if available)
```

Check `multipass find` for current blueprint availability — it varies by Multipass version and platform.

## Architecture Notes

- Image availability may differ between **amd64** and **arm64** (Apple Silicon Macs)
- On arm64, some older images or blueprints may not be available
- The default image (latest LTS) is always available on both architectures
- When an image isn't available, `multipass find` won't list it

## Recommendations

- **Default choice**: `24.04` (latest LTS, best package support, most testing)
- **Need older packages**: `22.04` for software that requires older Ubuntu
- **Testing Ubuntu upgrades**: Use daily builds of the next release
- **Minimal footprint**: Core images, but note they have limited package availability
- **Don't specify an image** unless the user has a reason — Multipass defaults to latest LTS
