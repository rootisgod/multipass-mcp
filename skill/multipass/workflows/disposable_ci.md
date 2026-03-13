# Workflow: Disposable CI Environment

Step-by-step playbook for running tests in a fresh VM and cleaning up afterward. The goal is to leave no trace — launch, test, collect results, destroy.

## Steps

### 1. Determine What to Run

Clarify with the user:
- What test command or script to run?
- What language/runtime is needed? (pick cloud-init template or install manually)
- Where is the project source? (mount from host or transfer files)
- What artifacts to collect? (logs, coverage reports, test output)

### 2. Launch a Fresh VM

Use a disposable name with a timestamp or run ID:
- Name pattern: `ci-run-<timestamp>` or `ci-<project>-<short-hash>`
- Keep resources modest: 2 CPUs, 2G RAM, 10G disk (unless tests are resource-heavy)
- Use cloud-init if the test environment needs specific tools

```
# Example: MCP
multipass_launch with name="ci-run-20260313", cpus=2, memory="2G", disk="10G"
```

### 3. Wait Until Running

- MCP: `multipass_wait_until_running`
- CLI: Poll until Running state

If cloud-init was used, also wait for it to complete:
- Execute: `["cloud-init", "status", "--wait"]`

### 4. Get Source Code Into the VM

Two options:
- **Mount** (faster for large projects): Mount the project directory read-only
- **Transfer** (cleaner isolation): Transfer a tarball or specific files

For mount:
- MCP: `multipass_mount_directory`
- CLI: `multipass mount /path/to/project ci-run-20260313:/home/ubuntu/project`

### 5. Run the Tests

Execute the test suite inside the VM:
- MCP: `multipass_run_script` with the test commands as a script, or `multipass_exec_command` for single commands
- CLI: `multipass exec ci-run-20260313 -- bash -c "cd /home/ubuntu/project && npm test"`

Set an appropriate timeout — test suites can take minutes.

Capture the output. If using `multipass_run_script`, the output is returned directly. If using CLI, redirect to a file and transfer it back.

### 6. Collect Artifacts

Transfer any test artifacts back to the host:
- MCP: `multipass_transfer` with source=`ci-run-20260313:/path/to/artifact` and destination=`/host/path/`
- CLI: `multipass transfer ci-run-20260313:/home/ubuntu/project/coverage/ ./coverage/ --recursive`

Common artifacts:
- Test output / JUnit XML
- Coverage reports
- Build logs
- Error logs

### 7. Tear Down

Unmount first (if mounted), then delete with purge:

1. Unmount:
   - MCP: `multipass_umount_directory`
   - CLI: `multipass umount ci-run-20260313:/home/ubuntu/project`

2. Delete and purge:
   - MCP: `multipass_delete` with `name="ci-run-20260313"` and `purge=true`
   - CLI: `multipass delete --purge ci-run-20260313`

### 8. Report Results

Present to the user:
- Test pass/fail summary
- Any collected artifacts and where they were saved
- Confirmation that the VM was deleted

## Key Principles

- **No trace**: The VM should be completely gone after the run
- **Reproducible**: Same inputs should produce same results in a fresh VM
- **Fast feedback**: Don't over-provision — use minimal resources for quick startup
- **Capture everything**: Transfer all artifacts before deleting — once purged, they're gone forever
