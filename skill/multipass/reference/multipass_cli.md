# Multipass CLI Reference

Cheat sheet for the `multipass` command. Use `--format json` on most commands for structured output.

## Instance Lifecycle

### launch
Create and start a new VM.
```
multipass launch [image] --name dev-vm --cpus 2 --memory 4G --disk 20G
multipass launch 24.04 --name docker-vm --cloud-init cloud-init.yaml
multipass launch --name mounted-vm --mount /host/path:/vm/path
multipass launch --name bridged-vm --network bridged
```
Key flags: `--name`, `--cpus`, `--memory`, `--disk`, `--cloud-init <file|url>`, `--network <name>`, `--mount <src:target>`, `--timeout <seconds>` (default 300)

### start / stop / restart / suspend
```
multipass start my-vm          # Start one instance
multipass start --all          # Start all instances
multipass stop my-vm           # Graceful stop
multipass stop --force my-vm   # Force stop
multipass restart my-vm
multipass suspend my-vm        # Save state to disk
```

### delete / recover / purge
```
multipass delete my-vm         # Move to trash (recoverable)
multipass delete --purge my-vm # Delete permanently
multipass delete --all         # Trash all instances
multipass recover my-vm        # Recover from trash
multipass recover --all        # Recover all trashed instances
multipass purge                # Permanently remove ALL trashed instances
```

## Command Execution

### exec
```
multipass exec my-vm -- ls -la /tmp
multipass exec my-vm -- bash -c "echo hello | grep hello"
multipass exec my-vm --working-directory /opt -- ./run.sh
```
Everything after `--` is the command. Use `bash -c "..."` for pipes/redirects.

### shell
```
multipass shell my-vm          # Interactive shell (not useful for automation)
```

## File Operations

### transfer
```
multipass transfer file.txt my-vm:/home/ubuntu/       # Host → VM
multipass transfer my-vm:/home/ubuntu/out.log ./       # VM → Host
multipass transfer --recursive ./dir my-vm:/home/ubuntu/dir  # Directories
```
Use `instance:path` syntax. Supports stdin/stdout with `-`.

### mount / umount
```
multipass mount /host/project my-vm:/home/ubuntu/project
multipass mount /data my-vm:/mnt/data --uid-map 1000:0 --gid-map 1000:0
multipass mount /src my-vm:/src --type native
multipass umount my-vm:/home/ubuntu/project
```
Key flags: `--uid-map <host:instance>`, `--gid-map <host:instance>`, `--type <classic|native>`

If mount fails with authorization error: `multipass set local.privileged-mounts=true`

## Snapshots

### snapshot / restore / clone
```
multipass snapshot my-vm                              # Auto-named snapshot
multipass snapshot my-vm --name before-docker --comment "Clean state"
multipass restore my-vm.before-docker                 # Restore (instance must be stopped)
multipass restore my-vm.before-docker --destructive   # Discard current state
multipass clone my-vm                                 # Clone (source must be stopped)
multipass clone my-vm --name my-vm-copy
```
Instance **must be stopped** for snapshot, restore, and clone.

## Information

### list / info / find
```
multipass list                           # All instances
multipass list --format json             # JSON output
multipass info my-vm                     # Detailed instance info
multipass info my-vm --format json       # JSON detailed info
multipass info my-vm --snapshots         # List snapshots
multipass find                           # Available images
multipass networks                       # Host network interfaces
```

### version / aliases
```
multipass version                        # Version info
multipass aliases                        # Command aliases
```

## Configuration

### get / set
```
multipass get --keys                     # List all setting keys
multipass get local.driver               # Read a specific setting
multipass set local.bridged-network=en0  # Set a value
multipass set local.privileged-mounts=true  # Enable privileged mounts
```

Common keys: `local.driver`, `local.bridged-network`, `local.privileged-mounts`

## authenticate
```
multipass authenticate <passphrase>      # Authenticate with daemon
```
