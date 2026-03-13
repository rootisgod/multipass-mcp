# Multipass MCP Tools Reference

The **multipass-mcp** server exposes 32 tools and 6 resources. When these tools are available, prefer them over raw CLI commands — they return structured JSON and have safety annotations.

- **Repository**: https://github.com/rootisgod/multipass-mcp
- **Install via Homebrew**: `brew install rootisgod/tap/multipass-mcp`
- **Register with Claude Code**: `claude mcp add multipass-mcp -- multipass-mcp`

## How to Detect

Check if `multipass_list_instances` is in your available tools. If yes, the MCP server is connected. If not, suggest the user install it using the commands above.

## Destructive Tools (confirm before using)

These tools have `destructiveHint: true` — always confirm with the user:

- `multipass_delete` — Deletes instances (recoverable unless `purge=true`)
- `multipass_restore` — Restores snapshot, discards current state when `destructive=true`
- `multipass_purge` — Permanently removes ALL trashed instances (irreversible)

---

## Instance Lifecycle (7 tools)

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_launch` | Launch a new VM | `image`, `name`, `cpus` (number), `disk`, `memory`, `cloud_init` (path/URL), `network`, `mount` (source:target), `timeout` (default 600s) |
| `multipass_start` | Start a stopped instance | `name` (required, or "all") |
| `multipass_stop` | Stop a running instance | `name` (required, or "all"), `force` (bool) |
| `multipass_restart` | Restart an instance | `name` (required, or "all") |
| `multipass_suspend` | Suspend a running instance | `name` (required, or "all") |
| `multipass_delete` | Delete an instance | `name` (required, or "all"), `purge` (bool) — **DESTRUCTIVE** |
| `multipass_recover` | Recover a deleted instance | `name` (required, or "all") |

Note: Pass `name="all"` on lifecycle tools to target all instances (translates to `--all` flag).

## Command Execution (2 tools)

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_exec_command` | Execute a command in a VM | `name` (required), `command` (required, string array e.g. `["ls", "-la"]`), `working_directory` |
| `multipass_run_script` | Run a multi-line script in a VM | `name` (required), `script` (required, string), `interpreter` (default "bash"), `working_directory`, `timeout` (default 300s) |

For complex shell commands with pipes/redirects, use: `command: ["bash", "-c", "your | complex command"]`

## File Operations (3 tools)

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_transfer` | Transfer files host↔VM | `source` (required), `destination` (required) — use `instance:/path` syntax, `recursive` (bool) |
| `multipass_mount_directory` | Mount host directory in VM | `source` (required, host path), `target` (required, `instance:/path`), `uid_map`, `gid_map`, `mount_type` ("classic"/"native") |
| `multipass_umount_directory` | Unmount a directory | `mount_path` (required, `instance:/path` or host path) |

## Snapshots (3 tools)

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_snapshot` | Create a snapshot (instance must be stopped) | `instance` (required), `name`, `comment` |
| `multipass_restore` | Restore to a snapshot (instance must be stopped) | `instance` (required), `snapshot` (required), `destructive` (bool) — **DESTRUCTIVE** |
| `multipass_clone` | Clone an instance (must be stopped) | `source_name` (required), `name` |

## Info & Queries (13 tools)

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_list_instances` | List all instances with state, IPs, image | (none) |
| `multipass_get_instance` | Detailed info for one instance | `name` (required) |
| `multipass_find_images` | List available images | (none) |
| `multipass_list_networks` | List host network interfaces | (none) |
| `multipass_list_snapshots` | List snapshots for an instance | `name` (required) |
| `multipass_get_version` | Get Multipass version info | (none) |
| `multipass_list_aliases` | List command aliases | (none) |
| `multipass_list_mounts` | List active mounts for an instance | `name` (required) |
| `multipass_list_deleted` | List deleted (trashed) instances | (none) |
| `multipass_instance_exists` | Check if an instance exists | `name` (required) — returns `{exists, state, image}` |
| `multipass_get_bridged_network` | Get configured bridged network | (none) |
| `multipass_disk_usage_check` | Check disk usage with threshold | `name` (required), `warn_percent` (default 80) |
| `multipass_wait_until_running` | Poll until instance is Running | `name` (required), `timeout` (default 120s) |

## Configuration (2 tools)

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_get_config` | Get a config value or list all keys | `key` (optional — omit to list all keys) |
| `multipass_set_config` | Set a config value | `key` (required), `value` (required) |

## System (2 tools)

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_purge` | Permanently delete all trashed instances | (none) — **DESTRUCTIVE, IRREVERSIBLE** |
| `multipass_authenticate` | Authenticate with Multipass service | `passphrase` (required) |

---

## Resources (6)

Resources provide read-only data access:

| URI | Description |
|-----|-------------|
| `multipass://instances` | All instances with state, IPs, image, resources |
| `multipass://images` | Available launchable images |
| `multipass://networks` | Host network interfaces |
| `multipass://version` | Multipass version info |
| `multipass://aliases` | Command aliases |
| `multipass://instance/{name}` | Detailed info for a specific instance (template) |

All resources return `application/json`.
