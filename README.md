# multipass-mcp

An MCP server for managing [Multipass](https://multipass.run/) virtual machines through AI assistants. Provides broad coverage of the Multipass CLI as MCP tools and resources, letting Claude (or any MCP client) launch, configure, snapshot, and manage Ubuntu VMs.

> **Note:** A few interactive/niche CLI verbs (`shell`, `alias`, `unalias`, `prefer`) are not exposed as tools since they either require a TTY or are rarely needed by AI assistants.

## Prerequisites

- [Multipass](https://multipass.run/) installed and on your PATH
- Verify with: `multipass version`

## Installation

### Homebrew (macOS/Linux)

```bash
brew install rootisgod/tap/multipass-mcp
```

### Go install

```bash
go install github.com/rootisgod/multipass-mcp@latest
```

### Download binary

Download the latest binary for your platform from the [releases page](https://github.com/rootisgod/multipass-mcp/releases).

### Build from source

```bash
git clone https://github.com/rootisgod/multipass-mcp.git
cd multipass-mcp
CGO_ENABLED=0 go build -ldflags "-X main.version=$(git describe --tags --always)" -o multipass-mcp .
```

## Configuration

### Claude Code

```bash
claude mcp add multipass-mcp -- /path/to/multipass-mcp
```

### Claude Desktop

Add to your config file:

- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "multipass": {
      "command": "/path/to/multipass-mcp"
    }
  }
}
```

## Claude Skill (Optional)

A companion **Claude skill** is included in this repo at [`skill/multipass/`](skill/multipass/). The MCP server gives Claude the *ability* to manage VMs; the skill gives it the *judgment* of when and how — decision trees for resource sizing, cloud-init templates, workflow playbooks, and error recovery patterns.

### Quick Start

**1. Install Multipass** (if you haven't already):
```bash
# macOS
brew install --cask multipass

# Linux (snap)
sudo snap install multipass
```

**2. Install the MCP server:**
```bash
brew install rootisgod/tap/multipass-mcp
```

**3. Register the MCP server with Claude Code:**
```bash
claude mcp add multipass-mcp -- multipass-mcp
```

**4. Install the skill** (available to all your projects):
```bash
git clone https://github.com/rootisgod/multipass-mcp.git /tmp/multipass-mcp
cp -r /tmp/multipass-mcp/skill/multipass ~/.claude/skills/multipass
rm -rf /tmp/multipass-mcp
```

Or for a single project only, copy into your repo's `.claude/skills/` directory instead.

**5. Use it:**

Just ask Claude naturally — the skill triggers automatically:

- *"Make me a VM with Docker installed"*
- *"Set up a dev environment for my Node project"*
- *"Snapshot my VM before I try something risky"*
- *"Launch 3 VMs for testing a distributed system"*
- *"Run my tests in a clean Ubuntu and then destroy it"*

### What the Skill Includes

| Directory | Contents |
|-----------|----------|
| `reference/` | MCP tools reference, CLI cheat sheet, cloud-init guide, networking, image catalog |
| `templates/` | Cloud-init configs for Docker, Kubernetes (microk8s), devtools (Node/Python), LAMP, Python data science |
| `workflows/` | Step-by-step playbooks for dev environments, snapshot/restore, multi-VM clusters, disposable CI |
| `scripts/` | Health check, SSH key setup, cleanup utilities |
| `evals/` | 8 test scenarios for validating skill behavior |

The skill works with or without the MCP server — if the `multipass_*` tools aren't detected, it falls back to `multipass` CLI commands via bash.

## Resources

Resources provide read-only access to Multipass state. MCP clients can read these URIs to get current information without side effects.

| URI | Description |
|-----|-------------|
| `multipass://instances` | All instances with state, IPs, image, and resource usage |
| `multipass://instance/{name}` | Detailed info for a specific instance (replace `{name}`) |
| `multipass://images` | Available images that can be launched |
| `multipass://networks` | Host network devices available for bridging |
| `multipass://version` | Multipass version information |
| `multipass://aliases` | Configured command aliases |

## Tools

All tool names are prefixed with `multipass_` to avoid conflicts when used alongside other MCP servers.

### Instance Lifecycle

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_launch` | Create and start a new VM | `image`, `name`, `cpus`, `disk`, `memory`, `cloud_init`, `network`, `mount` |
| `multipass_start` | Start a stopped instance | `name` (or `"all"`) |
| `multipass_stop` | Stop a running instance | `name` (or `"all"`), `force` |
| `multipass_restart` | Restart an instance | `name` (or `"all"`) |
| `multipass_suspend` | Suspend an instance to disk | `name` (or `"all"`) |
| `multipass_delete` | Move an instance to trash | `name` (or `"all"`), `purge` |
| `multipass_recover` | Restore a trashed instance | `name` (or `"all"`) |

For `multipass_start`, `multipass_stop`, `multipass_restart`, `multipass_suspend`, `multipass_delete`, and `multipass_recover`, pass `name="all"` to operate on every instance.

### Command Execution

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_exec_command` | Run a command inside an instance | `name`, `command` (list of strings), `working_directory`, `timeout` |
| `multipass_run_script` | Run a multi-line script inside an instance | `name`, `script`, `interpreter`, `working_directory` |

The `command` parameter takes a list of strings (e.g. `["ls", "-la", "/tmp"]`) rather than a shell string. This avoids shell injection and ensures each argument is passed correctly.

### File Operations

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_transfer` | Copy files between host and instance | `source`, `destination`, `recursive` |
| `multipass_mount_directory` | Mount a host directory inside an instance | `source`, `target`, `uid_map`, `gid_map`, `mount_type` |
| `multipass_umount_directory` | Unmount a mounted directory | `mount_path` |

File paths use `<instance>:<path>` syntax for instance-side paths:
- Host to instance: `source="/tmp/file.txt"`, `destination="my-vm:/home/ubuntu/file.txt"`
- Instance to host: `source="my-vm:/var/log/syslog"`, `destination="/tmp/syslog"`

### Snapshots

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_snapshot` | Create a snapshot of an instance | `instance`, `name`, `comment` |
| `multipass_restore` | Restore an instance to a snapshot | `instance`, `snapshot`, `destructive` |
| `multipass_clone` | Create an independent copy of an instance | `source_name`, `name` |

> **Important:** The instance must be **stopped** before taking a snapshot, restoring, or cloning. The server will return an error if the instance is running.

### Information & Queries

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_list_instances` | List all instances with state and IPs | *(none)* |
| `multipass_get_instance` | Get detailed info for one instance | `name` |
| `multipass_find_images` | List available images to launch | *(none)* |
| `multipass_list_networks` | List host network interfaces | *(none)* |
| `multipass_list_snapshots` | List snapshots for an instance | `name` |
| `multipass_get_version` | Get Multipass version info | *(none)* |
| `multipass_list_aliases` | List command aliases | *(none)* |
| `multipass_list_mounts` | List active mounts for an instance | `name` |
| `multipass_list_deleted` | List trashed instances | *(none)* |
| `multipass_instance_exists` | Check if an instance exists | `name` |
| `multipass_get_bridged_network` | Get the bridged network interface | *(none)* |
| `multipass_disk_usage_check` | Check disk usage with threshold | `name`, `warn_percent` |
| `multipass_wait_until_running` | Poll until instance is running | `name`, `timeout` |

### Configuration

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_get_config` | Read Multipass settings | `key` (omit for all keys) |
| `multipass_set_config` | Change a setting | `key`, `value` |

Common config keys:
- `local.driver` — VM backend (`qemu`, `virtualbox`, etc.)
- `local.privileged-mounts` — allow privileged mounts (`true`/`false`)
- `client.primary-name` — default instance name

### System

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `multipass_purge` | Permanently delete all trashed instances | *(none)* |
| `multipass_authenticate` | Authenticate with the Multipass daemon | `passphrase` |

## Tool Annotations

All tools include MCP tool annotations to help clients understand their behavior:

| Annotation | Meaning |
|------------|---------|
| `readOnlyHint: true` | Tool only reads data (all info/query tools) |
| `destructiveHint: true` | Tool may cause data loss (`multipass_delete`, `multipass_purge`, `multipass_restore`) |
| `idempotentHint: true` | Repeated calls have no additional effect (e.g. `multipass_start` on a running instance) |
| `openWorldHint: false` | Tool only interacts with local Multipass, not external services |

## Example Workflows

### Launch a dev VM and run commands

```
You: Launch an Ubuntu 24.04 VM named "dev" with 4 CPUs, 8G memory, and 20G disk

Claude: [calls multipass_launch(image="24.04", name="dev", cpus=4, memory="8G", disk="20G")]

You: Install Node.js in it

Claude: [calls multipass_exec_command(name="dev", command=["sudo", "apt-get", "update"])]
       [calls multipass_exec_command(name="dev", command=["sudo", "apt-get", "install", "-y", "nodejs", "npm"])]
```

### Snapshot and restore workflow

```
You: Snapshot the dev VM before I make risky changes

Claude: [calls multipass_stop(name="dev")]
       [calls multipass_snapshot(instance="dev", name="before-changes", comment="Clean state before experiment")]

You: Something broke, roll it back

Claude: [calls multipass_stop(name="dev")]
       [calls multipass_restore(instance="dev", snapshot="before-changes", destructive=True)]
       [calls multipass_start(name="dev")]
```

### Transfer files

```
You: Copy my local config to the VM

Claude: [calls multipass_transfer(source="/Users/me/.config/app/config.yaml",
                                  destination="dev:/home/ubuntu/.config/app/config.yaml")]
```

### Mount a project directory

```
You: Mount my project folder into the VM

Claude: [calls multipass_mount_directory(source="/Users/me/projects/myapp",
                                         target="dev:/home/ubuntu/myapp")]
```

### Clone a template VM

```
You: Clone the dev VM so I have a clean copy for testing

Claude: [calls multipass_stop(name="dev")]
       [calls multipass_clone(source_name="dev", name="dev-test")]
       [calls multipass_start(name="dev")]
       [calls multipass_start(name="dev-test")]
```

## Design Notes

- **No shell injection** — `multipass_exec_command` takes an array of strings (e.g. `["ls", "-la"]`), not a shell string. Arguments are passed directly to the subprocess.
- **Timeouts** — Most commands time out after 300s. `multipass_launch` defaults to 600s since image downloads can be slow.
- **Error handling** — CLI errors are returned with the stderr output, which MCP surfaces to the AI assistant automatically.
- **Stdio transport** — The server communicates over stdin/stdout using the MCP stdio protocol. Nothing is printed to stdout except MCP messages.
- **Static binary** — Built with `CGO_ENABLED=0` for a self-contained binary with no runtime dependencies.
- **Version injection** — The version is set at build time via `-ldflags "-X main.version=..."`. Defaults to `dev` for local builds.
