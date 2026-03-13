# Workflow: Snapshot and Restore

Step-by-step playbook for safe experimentation using snapshots.

## Creating a Snapshot (Before Risky Changes)

### 1. Confirm Target Instance

Verify which instance to snapshot:
- MCP: `multipass_instance_exists` with the name
- CLI: `multipass list`

If the instance doesn't exist, stop and clarify with the user.

### 2. Stop the Instance

Snapshots require the instance to be stopped:
- MCP: `multipass_stop` with the instance name
- CLI: `multipass stop <name>`

Wait for the stop to complete. If it hangs, suggest `force=true`.

### 3. Create the Snapshot

Use a descriptive name and comment:
- MCP: `multipass_snapshot` with `instance`, `name` (e.g., "before-docker-install"), and `comment` (e.g., "Clean state before adding Docker")
- CLI: `multipass snapshot <instance> --name before-docker-install --comment "Clean state before adding Docker"`

Good snapshot names describe what state they capture: `pre-migration`, `working-api`, `clean-slate`, `before-experiment`.

### 4. Restart the Instance

- MCP: `multipass_start` with the instance name
- CLI: `multipass start <name>`

### 5. Confirm to User

Report: snapshot name, instance name, and that they can now proceed with their risky changes.

---

## Restoring a Snapshot (Rolling Back)

### 1. Identify the Snapshot

List available snapshots:
- MCP: `multipass_list_snapshots` with the instance name
- CLI: `multipass info <name> --snapshots`

If there are multiple snapshots, present the list and ask which one to restore.

### 2. Stop the Instance

- MCP: `multipass_stop` with the instance name
- CLI: `multipass stop <name>`

### 3. Restore the Snapshot

This is a destructive operation — the current state is discarded:
- MCP: `multipass_restore` with `instance`, `snapshot`, and `destructive=true`
- CLI: `multipass restore <instance>.<snapshot> --destructive`

**Always confirm with the user** before restoring, since it discards all changes since the snapshot.

### 4. Restart the Instance

- MCP: `multipass_start` with the instance name
- CLI: `multipass start <name>`

### 5. Verify

Run a quick health check or verify the specific state the user expected.

Report: which snapshot was restored, and that the instance is running.

---

## Pre-Flight Check Pattern

Use this pattern whenever doing something potentially destructive in a VM:

```
1. Stop instance
2. Snapshot with descriptive name
3. Start instance
4. Do the risky thing
5. If broken → stop → restore --destructive → start
6. If good → continue (snapshot remains as rollback point)
```

This is cheap insurance. Snapshots take seconds and can save hours of rebuilding.
