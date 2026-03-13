# Workflow: Set Up a Development Environment

Step-by-step playbook for creating a dev VM tailored to a project.

## Steps

### 1. Gather Requirements

Ask the user (if not already clear):
- What project or language? (Node, Python, Go, Rust, general, etc.)
- Do they have a local project directory to mount?
- Any specific tools or services needed? (Docker, databases, etc.)

### 2. Choose Configuration

Based on the answers:
- **Image**: `24.04` (latest LTS) unless they need an older Ubuntu
- **Resources**: 4 CPUs, 8G RAM, 20G disk for dev work. Scale down for lightweight tasks.
- **Cloud-init**: Check if a template from `templates/` fits their stack:
  - Node/JavaScript → `cloud-init-devtools.yaml`
  - Docker → `cloud-init-docker.yaml`
  - Kubernetes → `cloud-init-k8s.yaml`
  - PHP/web → `cloud-init-lamp.yaml`
  - Python data science → `cloud-init-python-ds.yaml`
  - If no template fits, prepare a custom cloud-init or install tools via exec after launch

### 3. Check for Name Conflicts

Before launching, verify the chosen name doesn't already exist:
- MCP: `multipass_instance_exists` with the intended name
- CLI: `multipass list` and check output

If it exists, ask the user: use the existing one, pick a different name, or delete and recreate?

### 4. Launch the VM

Use a descriptive name like `dev-<project>`:
- MCP: `multipass_launch` with `name`, `cpus`, `memory`, `disk`, and optionally `cloud_init`
- CLI: `multipass launch 24.04 --name dev-myproject --cpus 4 --memory 8G --disk 20G`

### 5. Wait Until Running

- MCP: `multipass_wait_until_running` with the instance name
- CLI: Poll `multipass list` until state is "Running"

### 6. Verify Health

Transfer and run the health check script:
1. Read `scripts/health_check.sh`
2. Transfer it to the VM
3. Execute it and check the JSON output
4. Confirm: apt works, DNS resolves, sufficient disk and memory

If cloud-init was used, also verify it completed:
- Execute: `["cloud-init", "status", "--wait"]`
- Check marker: `["test", "-f", "/var/run/cloud-init-complete"]`

### 7. Mount Project Directory

If the user has a local project directory:
- MCP: `multipass_mount_directory` with `source` (host path) and `target` (`instance:/home/ubuntu/project`)
- CLI: `multipass mount /path/to/project dev-myproject:/home/ubuntu/project`

If mount fails, check `local.privileged-mounts` and suggest enabling it.

### 8. Install Additional Tools

If the cloud-init didn't cover everything, install remaining tools via exec:
- MCP: `multipass_exec_command` or `multipass_run_script`
- CLI: `multipass exec dev-myproject -- sudo apt-get install -y <packages>`

Verify each installation (`which <tool>`, `<tool> --version`).

### 9. Report to User

Provide:
- Instance name
- IP address (from list_instances or info)
- What was installed
- Mount point (if applicable)
- How to connect: `multipass shell <name>` or `ssh ubuntu@<ip>`

## Checklist

- [ ] Requirements gathered (language, tools, project dir)
- [ ] Configuration chosen (image, resources, cloud-init)
- [ ] Name conflict checked
- [ ] VM launched successfully
- [ ] VM is running and healthy
- [ ] Cloud-init completed (if used)
- [ ] Project directory mounted (if requested)
- [ ] Additional tools installed and verified
- [ ] Connection info reported to user
