---
name: multipass
description: >
  Manage Canonical Multipass Ubuntu virtual machines — launch, configure,
  snapshot, and orchestrate VMs for development, testing, and CI workflows.
  Use this skill whenever the user mentions Multipass, Ubuntu VMs, virtual
  machines for development, test environments, cloud-init, disposable Linux
  instances, or wants to run commands in an isolated Ubuntu environment.
  Also trigger when the user asks to "set up a dev VM", "create a test
  environment", "spin up a Linux box", "I need an isolated environment",
  "run this in a clean Ubuntu", or similar — even if they don't mention
  Multipass by name. Trigger for anything involving VM lifecycle (launch,
  stop, snapshot, clone), cloud-init configuration, or multi-VM orchestration.
---

# Multipass Skill

Manage Canonical Multipass Ubuntu VMs — from single dev environments to multi-node clusters. This skill works with or without the `multipass-mcp` MCP server. When the MCP server is connected, its 32 tools provide structured JSON responses with safety annotations. Without it, fall back to `multipass` CLI commands via bash. Either way, the workflows and judgment patterns in this skill apply.

## Companion MCP Server

This skill is designed to work alongside the **multipass-mcp** MCP server, which exposes the full Multipass CLI as 32 structured tools with safety annotations.

- **Repository**: https://github.com/rootisgod/multipass-mcp
- **Install via Homebrew**: `brew install rootisgod/tap/multipass-mcp`
- **Register with Claude Code**: `claude mcp add multipass-mcp -- multipass-mcp`

The MCP server is optional but strongly recommended — it gives you structured JSON responses, proper error handling, and destructive-action annotations that raw CLI calls lack. If the MCP tools aren't detected, the skill falls back to bash `multipass` commands automatically.

## MCP Detection

At the start of any Multipass workflow, determine which interface to use:

```
1. Check: Is the multipass_list_instances tool available?
   → YES: Use multipass_* MCP tools throughout (preferred — structured JSON, error handling)
   → NO:  Use bash with `multipass` CLI commands (append --format json where useful)
```

If MCP tools are available, read `reference/mcp_tools.md` to understand tool names, parameters, and which ones are destructive. If not, read `reference/multipass_cli.md` for CLI flag reference.

## Decision Tree — What to Launch

Match user intent to the right configuration:

| Intent | Image | CPUs | RAM | Disk | Extras |
|--------|-------|------|-----|------|--------|
| Quick test / throwaway | 24.04 LTS | 1 | 1G | 5G | No mounts, plan to delete after |
| General dev environment | 24.04 LTS | 2 | 4G | 20G | Mount project dir |
| Heavy dev (compilation, Docker) | 24.04 LTS | 4 | 8G | 30G | Mount project dir, cloud-init |
| Specific stack (Docker/K8s/LAMP) | 24.04 LTS | 2-4 | 4-8G | 20G | Use matching cloud-init template |
| Multiple VMs / cluster | 24.04 LTS | 2 each | 2G each | 10G each | Sequential naming |

When the user doesn't specify resources, default to **2 CPUs, 2G memory, 10G disk** for light tasks and **4 CPUs, 8G memory, 20G disk** for development work.

Read `reference/image_catalog.md` when choosing an image. Read `reference/cloud_init_guide.md` plus the relevant template from `templates/` when using cloud-init.

## Core Patterns

### Always Snapshot Before Risky Changes

Before installing unfamiliar software, running untested scripts, or modifying system configs — stop the VM and snapshot it. Snapshots are cheap insurance against irreversible mistakes. The pattern is: stop → snapshot → start → do risky thing → if broken: stop → restore --destructive → start.

### Verify After Action

Don't assume commands succeeded. After launching, verify the VM is running and can reach the internet. After installing software, check it's actually installed (`which docker`, `node --version`, etc.). Transfer and run `scripts/health_check.sh` for a comprehensive check when setting up a new environment.

### Clean Naming

Name instances descriptively based on their purpose: `dev-myproject`, `test-api`, `ci-run-42`, `docker-lab`. Avoid generic names like `primary`, `ubuntu`, `test1`. Check if the name already exists with `multipass_instance_exists` (MCP) or `multipass list` (CLI) before launching.

### Mount Over Transfer

When the user is actively developing, mount their project directory into the VM rather than transferring files. Mounts stay in sync automatically and are faster for iterative work. Use transfer only for one-off file moves or artifacts.

### Confirm Before Destructive Actions

Tools marked `destructiveHint: true` (delete, restore, purge) can cause data loss. Always confirm with the user before using them. When deleting, prefer `delete` without `purge` first — instances can be recovered. Only purge when the user explicitly wants permanent removal.

## Reference Loading

Load reference files on demand based on what you're doing:

| Situation | Load |
|-----------|------|
| Starting any Multipass workflow | `reference/mcp_tools.md` (detect MCP availability) |
| Choosing an image to launch | `reference/image_catalog.md` |
| Using cloud-init | `reference/cloud_init_guide.md` + relevant `templates/*.yaml` |
| Networking issues or bridged setup | `reference/networking.md` |
| Looking up CLI flags | `reference/multipass_cli.md` |

## Workflow Dispatch

Map user intents to workflow playbooks:

| User says something like... | Workflow |
|-----------------------------|----------|
| "Set up a dev VM", "create a development environment" | `workflows/dev_environment.md` |
| "Snapshot my VM", "save state before changes", "backup" | `workflows/snapshot_restore.md` |
| "I need multiple VMs", "set up a cluster", "launch N nodes" | `workflows/multi_vm_cluster.md` |
| "Run tests in a clean environment", "disposable VM" | `workflows/disposable_ci.md` |

For requests that don't match a workflow (e.g., "list my VMs", "stop everything"), handle directly using the appropriate MCP tool or CLI command.

## Error Recovery

| Error | Cause | Fix |
|-------|-------|-----|
| "instance not found" | Typo or instance was deleted | Run `multipass_list_instances` to show what exists; suggest closest match |
| "timed out waiting for response" | Slow image download or overloaded host | Retry with longer timeout; check host resources |
| "cannot snapshot running instance" | Snapshot requires stopped state | Stop the instance first, then snapshot |
| "mount failed" / "not authorized" | Privileged mounts disabled | Run `multipass_set_config` with `key=local.privileged-mounts, value=true` |
| "not enough disk space" | Deleted instances consuming space | Run `multipass_purge` to reclaim space from trashed instances |
| "instance already exists" | Name collision | Pick a different name or check if the existing one is what they want |
| "launch failed" / network errors | DNS or connectivity issues inside new VM | Run health check script; check host network |

## Scripts

Three utility scripts are available in `scripts/`:

- **`health_check.sh`** — Transfer into a VM and run to verify apt, DNS, disk, and memory. Returns JSON summary. Use after launch or when diagnosing issues.
- **`setup_ssh_keys.sh`** — Transfer into a VM and run with a public key as argument. Sets up `~/.ssh/authorized_keys` with proper permissions.
- **`cleanup_all.sh`** — Run on the HOST (not in a VM). Stops all instances, deletes all, and purges. Use only when the user explicitly wants a clean slate.

To use a script: read it from `scripts/`, transfer it into the VM with `multipass_transfer` or `multipass transfer`, then execute it with `multipass_exec_command` or `multipass exec`.
