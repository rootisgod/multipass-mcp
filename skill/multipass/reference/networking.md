# Multipass Networking Reference

## Default NAT Networking

By default, Multipass VMs get a private IP on an internal NAT network:

- VM gets an IP (typically `192.168.64.x` on macOS, `10.x.x.x` on Linux)
- Host can reach VM via this IP
- VM can reach the internet
- Other machines on the LAN cannot reach the VM directly

To find a VM's IP:
```bash
multipass list                    # Shows IP in the IPv4 column
multipass info my-vm              # Detailed info including IP
```

With MCP: `multipass_list_instances` or `multipass_get_instance` returns IPs in JSON.

## Bridged Networking

Bridged networking gives VMs an IP on your host's LAN, making them accessible to other machines.

### Setup Steps

1. **List available networks** to find the right interface:
   ```bash
   multipass networks
   ```
   With MCP: `multipass_list_networks`

2. **Configure the bridged interface** (one-time):
   ```bash
   multipass set local.bridged-network=en0    # macOS example
   multipass set local.bridged-network=eth0   # Linux example
   ```
   With MCP: `multipass_set_config` with `key=local.bridged-network`, `value=en0`

3. **Launch with bridged networking**:
   ```bash
   multipass launch --name bridged-vm --network bridged
   multipass launch --name bridged-vm --network name=en0   # Explicit interface
   ```
   With MCP: `multipass_launch` with `network=bridged` or `network=name=en0`

### Checking Bridged Network Config

```bash
multipass get local.bridged-network
```

With MCP: `multipass_get_bridged_network` — returns the configured interface or null.

## SSH from Host to VM

Multipass VMs are already accessible via `multipass exec` and `multipass shell`, but if you need direct SSH:

1. The `ubuntu` user exists by default with passwordless sudo
2. Get the VM's IP from `multipass list`
3. Multipass automatically sets up SSH keys — you can SSH with:
   ```bash
   ssh ubuntu@<vm-ip>
   ```
4. To add custom SSH keys, use the `scripts/setup_ssh_keys.sh` script or cloud-init `ssh_authorized_keys`

## Port Forwarding

Multipass doesn't have built-in port forwarding. Options:

- **Access directly via VM IP**: From the host, connect to `<vm-ip>:<port>` (works with NAT)
- **Use bridged networking**: VM gets a LAN IP accessible from anywhere on the network
- **SSH tunnel**: `ssh -L 8080:localhost:80 ubuntu@<vm-ip>` forwards host:8080 → vm:80

## Common Gotchas

- `multipass networks` must show available interfaces before you can use `--network`
- On macOS, the interface is typically `en0` (Wi-Fi) or `en1`
- On Linux, common names are `eth0`, `enp0s3`, `wlp2s0`
- Bridged networking may not work on all network types (some corporate Wi-Fi blocks it)
- If a VM has no IP after launch, wait a moment — DHCP can take a few seconds
- VMs behind NAT can reach the internet but aren't reachable from other LAN hosts
