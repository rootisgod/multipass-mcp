# Cloud-Init Guide for Multipass

Cloud-init runs on first boot of a VM to automate initial setup. Multipass supports passing cloud-init config at launch time.

## Usage

```bash
multipass launch --name my-vm --cloud-init config.yaml
multipass launch --name my-vm --cloud-init https://example.com/config.yaml
```

With MCP: pass the file path as the `cloud_init` parameter to `multipass_launch`.

## Syntax

Files must start with `#cloud-config` on the first line. The rest is YAML.

### Install Packages

```yaml
#cloud-config
package_update: true
package_upgrade: true
packages:
  - git
  - curl
  - build-essential
```

### Run Commands on First Boot

```yaml
#cloud-config
runcmd:
  - echo "Setup started" > /var/log/setup.log
  - apt-get update
  - apt-get install -y nginx
  - systemctl enable nginx
  - touch /var/run/cloud-init-complete
```

Commands run as root. Each list item is either a string (run via shell) or a list (run directly).

### Write Files

```yaml
#cloud-config
write_files:
  - path: /etc/myapp/config.json
    content: |
      {"port": 8080, "debug": false}
    permissions: "0644"
    owner: root:root
  - path: /home/ubuntu/.bashrc
    content: |
      export PATH=$PATH:/opt/bin
    append: true
    owner: ubuntu:ubuntu
```

### Configure Users and SSH

```yaml
#cloud-config
users:
  - default
  - name: deploy
    groups: sudo, docker
    shell: /bin/bash
    ssh_authorized_keys:
      - ssh-rsa AAAA... user@host
```

The `default` entry preserves the default `ubuntu` user.

### Combined Example

```yaml
#cloud-config
package_update: true
packages:
  - git
  - docker.io

runcmd:
  - usermod -aG docker ubuntu
  - systemctl enable docker
  - systemctl start docker
  - touch /var/run/cloud-init-complete

write_files:
  - path: /home/ubuntu/.env
    content: |
      ENVIRONMENT=development
    owner: ubuntu:ubuntu
```

## Important Notes

- Cloud-init runs **once on first boot only**. Relaunching or restarting won't re-run it.
- Logs are at `/var/log/cloud-init-output.log` — check here if setup seems incomplete.
- Cloud-init can take several minutes for package installs. Wait for completion before verifying.
- To check if cloud-init finished: `cloud-init status --wait` (blocks until done) or check for a marker file like `/var/run/cloud-init-complete`.
- The `ubuntu` user is the default user in Multipass VMs and has passwordless sudo.

## Verification Pattern

After launching with cloud-init:
1. Wait for the instance to be running (`multipass_wait_until_running` or `multipass list`)
2. Check cloud-init status: `exec` → `["cloud-init", "status", "--wait"]`
3. Verify installed software: `exec` → `["which", "docker"]` or similar
4. Check the marker file: `exec` → `["test", "-f", "/var/run/cloud-init-complete"]`
