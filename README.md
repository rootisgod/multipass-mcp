# multipass-mcp

An MCP server for managing [Multipass](https://multipass.run/) virtual machines through AI assistants. Exposes the full Multipass CLI as MCP tools and resources, letting Claude (or any MCP client) launch, configure, snapshot, and manage Ubuntu VMs.

## Prerequisites

- [Multipass](https://multipass.run/) installed and on your PATH
- Verify with: `multipass version`

## Installation

### Homebrew (macOS/Linux)

```bash
brew install iainmckee/tap/multipass-mcp
```

### Download binary

Download the latest binary for your platform from the [releases page](https://github.com/iainmckee/multipass-mcp/releases).

### Build from source

```bash
git clone https://github.com/iainmckee/multipass-mcp.git
cd multipass-mcp
CGO_ENABLED=0 go build -o multipass-mcp .
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

### Instance Lifecycle

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `launch` | Create and start a new VM | `image`, `name`, `cpus`, `disk`, `memory`, `cloud_init`, `network`, `mount` |
| `start` | Start a stopped instance | `name` (or `"all"`) |
| `stop` | Stop a running instance | `name` (or `"all"`), `force` |
| `restart` | Restart an instance | `name` (or `"all"`) |
| `suspend` | Suspend an instance to disk | `name` (or `"all"`) |
| `delete` | Move an instance to trash | `name` (or `"all"`), `purge` |
| `recover` | Restore a trashed instance | `name` (or `"all"`) |

For `start`, `stop`, `restart`, `suspend`, `delete`, and `recover`, pass `name="all"` to operate on every instance.

### Command Execution

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `exec_command` | Run a command inside an instance | `name`, `command` (list of strings), `working_directory` |

The `command` parameter takes a list of strings (e.g. `["ls", "-la", "/tmp"]`) rather than a shell string. This avoids shell injection and ensures each argument is passed correctly.

### File Operations

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `transfer` | Copy files between host and instance | `source`, `destination`, `recursive` |
| `mount_directory` | Mount a host directory inside an instance | `source`, `target`, `uid_map`, `gid_map`, `mount_type` |
| `umount_directory` | Unmount a mounted directory | `mount_path` |

File paths use `<instance>:<path>` syntax for instance-side paths:
- Host to instance: `source="/tmp/file.txt"`, `destination="my-vm:/home/ubuntu/file.txt"`
- Instance to host: `source="my-vm:/var/log/syslog"`, `destination="/tmp/syslog"`

### Snapshots

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `snapshot` | Create a snapshot of an instance | `instance`, `name`, `comment` |
| `restore` | Restore an instance to a snapshot | `instance`, `snapshot`, `destructive` |
| `clone` | Create an independent copy of an instance | `source_name`, `name` |

> **Important:** The instance must be **stopped** before taking a snapshot, restoring, or cloning. The server will return an error if the instance is running.

### Configuration

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `get_config` | Read Multipass settings | `key` (omit for all keys) |
| `set_config` | Change a setting | `key`, `value` |

Common config keys:
- `local.driver` — VM backend (`qemu`, `virtualbox`, etc.)
- `local.privileged-mounts` — allow privileged mounts (`true`/`false`)
- `client.primary-name` — default instance name

### System

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `purge` | Permanently delete all trashed instances | *(none)* |
| `authenticate` | Authenticate with the Multipass daemon | `passphrase` |

## Example Workflows

### Launch a dev VM and run commands

```
You: Launch an Ubuntu 24.04 VM named "dev" with 4 CPUs, 8G memory, and 20G disk

Claude: [calls launch(image="24.04", name="dev", cpus=4, memory="8G", disk="20G")]

You: Install Node.js in it

Claude: [calls exec_command(name="dev", command=["sudo", "apt-get", "update"])]
       [calls exec_command(name="dev", command=["sudo", "apt-get", "install", "-y", "nodejs", "npm"])]
```

### Snapshot and restore workflow

```
You: Snapshot the dev VM before I make risky changes

Claude: [calls stop(name="dev")]
       [calls snapshot(instance="dev", name="before-changes", comment="Clean state before experiment")]

You: Something broke, roll it back

Claude: [calls stop(name="dev")]
       [calls restore(instance="dev", snapshot="before-changes", destructive=True)]
       [calls start(name="dev")]
```

### Transfer files

```
You: Copy my local config to the VM

Claude: [calls transfer(source="/Users/me/.config/app/config.yaml",
                        destination="dev:/home/ubuntu/.config/app/config.yaml")]
```

### Mount a project directory

```
You: Mount my project folder into the VM

Claude: [calls mount_directory(source="/Users/me/projects/myapp",
                               target="dev:/home/ubuntu/myapp")]
```

### Clone a template VM

```
You: Clone the dev VM so I have a clean copy for testing

Claude: [calls stop(name="dev")]
       [calls clone(source_name="dev", name="dev-test")]
       [calls start(name="dev")]
       [calls start(name="dev-test")]
```

## Design Notes

- **No shell injection** — `exec_command` takes an array of strings (e.g. `["ls", "-la"]`), not a shell string. Arguments are passed directly to the subprocess.
- **Timeouts** — Most commands time out after 300s. `launch` defaults to 600s since image downloads can be slow.
- **Error handling** — CLI errors are returned with the stderr output, which MCP surfaces to the AI assistant automatically.
- **Stdio transport** — The server communicates over stdin/stdout using the MCP stdio protocol. Nothing is printed to stdout except MCP messages.
- **Static binary** — Built with `CGO_ENABLED=0` for a self-contained binary with no runtime dependencies.
