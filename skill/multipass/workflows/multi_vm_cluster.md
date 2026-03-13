# Workflow: Multi-VM Cluster

Step-by-step playbook for launching and connecting multiple VMs.

## Steps

### 1. Gather Requirements

Ask the user (if not already clear):
- How many VMs?
- What are they for? (distributed system test, load balancing, microservices, etc.)
- Do they need to communicate with each other? (usually yes)
- Any specific software on each node?
- Do they need to be accessible from the host LAN? (bridged networking)

### 2. Plan Resources

Each VM consumes host resources. Be mindful of limits:
- **Per VM**: 2 CPUs, 2G RAM, 10G disk is a reasonable default for cluster nodes
- **Total check**: N VMs × resources should not exceed ~75% of host capacity
- For 3 nodes at 2 CPUs/2G each: 6 CPUs, 6G RAM needed — reasonable on most dev machines
- For 5+ nodes: consider reducing per-node resources

Warn the user if the total seems high for typical hardware.

### 3. Choose Naming Convention

Use sequential, descriptive names:
- `node-1`, `node-2`, `node-3` for generic clusters
- `web-1`, `web-2` and `db-1` for role-based naming
- `worker-1`, `worker-2`, `manager-1` for orchestration setups

### 4. Launch All VMs

Launch sequentially (Multipass doesn't support parallel launches well):

For each VM:
- MCP: `multipass_launch` with `name`, `cpus`, `memory`, `disk`
- CLI: `multipass launch 24.04 --name node-1 --cpus 2 --memory 2G --disk 10G`

If all nodes need the same software, use a cloud-init template to avoid post-launch setup.

### 5. Wait Until All Running

For each VM:
- MCP: `multipass_wait_until_running`
- CLI: Poll `multipass list` until all show "Running"

### 6. Collect IP Addresses

- MCP: `multipass_list_instances` — parse the JSON for all node IPs
- CLI: `multipass list --format json`

Record the IPs and present them to the user in a table:

```
| Name   | IP            |
|--------|---------------|
| node-1 | 192.168.64.2  |
| node-2 | 192.168.64.3  |
| node-3 | 192.168.64.4  |
```

### 7. Verify Inter-Node Connectivity

From node-1, ping the other nodes:
- MCP: `multipass_exec_command` on node-1 with `command: ["ping", "-c", "1", "<node-2-ip>"]`
- CLI: `multipass exec node-1 -- ping -c 1 <node-2-ip>`

If pings fail, check that all VMs are on the same network.

### 8. Optional: Configure Bridged Networking

If VMs need to be accessible from the host LAN:
1. Check available networks: `multipass_list_networks`
2. Set bridged interface: `multipass_set_config` with `key=local.bridged-network`
3. Relaunch with `--network bridged` (or add network to existing VMs — requires relaunch)

See `reference/networking.md` for details.

### 9. Optional: Install Software

If cloud-init wasn't used, install software on each node:
- For identical setups: write a script once, transfer and execute on each node
- For role-based setups: use different scripts per role

### 10. Report to User

Provide:
- Number of VMs launched
- IP address table
- Connectivity status (inter-node pings)
- Installed software (if any)
- Total host resources consumed

## Cleanup Note

When done, the user can tear down the cluster:
- MCP: `multipass_delete` with `name="all"` then `multipass_purge`
- CLI: `multipass delete --all && multipass purge`
- Or use `scripts/cleanup_all.sh`

Always confirm before deleting — the user may want to keep the cluster.
