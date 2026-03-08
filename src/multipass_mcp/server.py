"""MCP server for managing Multipass virtual machines."""

import asyncio
import json
import sys

from mcp.server.fastmcp import FastMCP

mcp = FastMCP("multipass")


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------


async def run_multipass(*args: str, timeout: int = 300) -> str:
    """Run a multipass command and return stdout. Raises RuntimeError on failure."""
    proc = await asyncio.create_subprocess_exec(
        "multipass",
        *args,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )
    try:
        stdout, stderr = await asyncio.wait_for(proc.communicate(), timeout=timeout)
    except asyncio.TimeoutError:
        proc.kill()
        await proc.communicate()
        raise RuntimeError(f"multipass {args[0]} timed out after {timeout}s")
    if proc.returncode != 0:
        raise RuntimeError(
            f"multipass {' '.join(args)} failed (exit {proc.returncode}): "
            f"{stderr.decode().strip()}"
        )
    return stdout.decode().strip()


async def run_multipass_json(*args: str, timeout: int = 300) -> dict | list:
    """Run a multipass command with --format json and return parsed result."""
    result = await run_multipass(*args, "--format", "json", timeout=timeout)
    return json.loads(result)


# ---------------------------------------------------------------------------
# Resources
# ---------------------------------------------------------------------------


@mcp.resource("multipass://instances")
async def list_instances() -> str:
    """List all Multipass instances with state, IPs, image, and resource usage."""
    data = await run_multipass_json("list")
    return json.dumps(data, indent=2)


@mcp.resource("multipass://instance/{name}")
async def get_instance(name: str) -> str:
    """Get detailed information about a specific Multipass instance."""
    data = await run_multipass_json("info", name)
    return json.dumps(data, indent=2)


@mcp.resource("multipass://images")
async def list_images() -> str:
    """List available Multipass images that can be launched."""
    data = await run_multipass_json("find")
    return json.dumps(data, indent=2)


@mcp.resource("multipass://networks")
async def list_networks() -> str:
    """List host network devices available for Multipass instances."""
    data = await run_multipass_json("networks")
    return json.dumps(data, indent=2)


@mcp.resource("multipass://version")
async def get_version() -> str:
    """Get Multipass version information."""
    data = await run_multipass_json("version")
    return json.dumps(data, indent=2)


@mcp.resource("multipass://aliases")
async def list_aliases() -> str:
    """List configured Multipass command aliases."""
    data = await run_multipass_json("aliases")
    return json.dumps(data, indent=2)


# ---------------------------------------------------------------------------
# Tools — Instance Lifecycle
# ---------------------------------------------------------------------------


@mcp.tool()
async def launch(
    image: str = "",
    name: str = "",
    cpus: int | None = None,
    disk: str = "",
    memory: str = "",
    cloud_init: str = "",
    network: str = "",
    mount: str = "",
    timeout: int = 600,
) -> str:
    """Launch a new Multipass VM instance.

    Args:
        image: Image to launch (e.g. "22.04", "daily:24.04"). Defaults to latest LTS.
        name: Name for the instance. Auto-generated if omitted.
        cpus: Number of CPUs to allocate.
        disk: Disk size (e.g. "10G", "50G").
        memory: Memory size (e.g. "1G", "4G").
        cloud_init: Path or URL to cloud-init config file.
        network: Network to connect to (from multipass networks).
        mount: Host path to mount in the instance (source:target format).
        timeout: Timeout in seconds (default 600 — image downloads can be slow).
    """
    args = ["launch"]
    if image:
        args.append(image)
    if name:
        args.extend(["--name", name])
    if cpus is not None:
        args.extend(["--cpus", str(cpus)])
    if disk:
        args.extend(["--disk", disk])
    if memory:
        args.extend(["--memory", memory])
    if cloud_init:
        args.extend(["--cloud-init", cloud_init])
    if network:
        args.extend(["--network", network])
    if mount:
        args.extend(["--mount", mount])
    return await run_multipass(*args, timeout=timeout)


@mcp.tool()
async def start(name: str) -> str:
    """Start a stopped Multipass instance.

    Args:
        name: Instance name, or "all" to start all instances.
    """
    if name == "all":
        return await run_multipass("start", "--all")
    return await run_multipass("start", name)


@mcp.tool()
async def stop(name: str, force: bool = False) -> str:
    """Stop a running Multipass instance.

    Args:
        name: Instance name, or "all" to stop all instances.
        force: Force stop without waiting for graceful shutdown.
    """
    args = ["stop"]
    if force:
        args.append("--force")
    if name == "all":
        args.append("--all")
    else:
        args.append(name)
    return await run_multipass(*args)


@mcp.tool()
async def restart(name: str) -> str:
    """Restart a Multipass instance.

    Args:
        name: Instance name, or "all" to restart all instances.
    """
    if name == "all":
        return await run_multipass("restart", "--all")
    return await run_multipass("restart", name)


@mcp.tool()
async def suspend(name: str) -> str:
    """Suspend a running Multipass instance.

    Args:
        name: Instance name, or "all" to suspend all instances.
    """
    if name == "all":
        return await run_multipass("suspend", "--all")
    return await run_multipass("suspend", name)


@mcp.tool()
async def delete(name: str, purge: bool = False) -> str:
    """Delete a Multipass instance (can be recovered unless purged).

    Args:
        name: Instance name, or "all" to delete all instances.
        purge: Permanently delete instead of moving to trash.
    """
    args = ["delete"]
    if purge:
        args.append("--purge")
    if name == "all":
        args.append("--all")
    else:
        args.append(name)
    return await run_multipass(*args)


@mcp.tool()
async def recover(name: str) -> str:
    """Recover a deleted (trashed) Multipass instance.

    Args:
        name: Instance name, or "all" to recover all deleted instances.
    """
    if name == "all":
        return await run_multipass("recover", "--all")
    return await run_multipass("recover", name)


# ---------------------------------------------------------------------------
# Tools — Execution
# ---------------------------------------------------------------------------


@mcp.tool()
async def exec_command(
    name: str,
    command: list[str],
    working_directory: str = "",
) -> str:
    """Execute a command inside a Multipass instance.

    Args:
        name: Instance name.
        command: Command and arguments as a list (e.g. ["ls", "-la", "/tmp"]).
        working_directory: Working directory inside the instance.
    """
    args = ["exec", name]
    if working_directory:
        args.extend(["--working-directory", working_directory])
    args.append("--")
    args.extend(command)
    return await run_multipass(*args)


# ---------------------------------------------------------------------------
# Tools — File Operations
# ---------------------------------------------------------------------------


@mcp.tool()
async def transfer(
    source: str,
    destination: str,
    recursive: bool = False,
) -> str:
    """Transfer files between host and a Multipass instance.

    Use <name>:<path> syntax for instance paths (e.g. "my-vm:/home/ubuntu/file.txt").
    Use a plain path for host paths (e.g. "/tmp/file.txt").

    Args:
        source: Source path (host path or instance:path).
        destination: Destination path (host path or instance:path).
        recursive: Transfer directories recursively.
    """
    args = ["transfer"]
    if recursive:
        args.append("--recursive")
    args.extend([source, destination])
    return await run_multipass(*args)


@mcp.tool()
async def mount_directory(
    source: str,
    target: str,
    uid_map: str = "",
    gid_map: str = "",
    mount_type: str = "",
) -> str:
    """Mount a host directory inside a Multipass instance.

    Args:
        source: Host directory path to mount.
        target: Mount point in instance:path format (e.g. "my-vm:/mnt/data").
        uid_map: UID mapping in host:instance format (e.g. "1000:0").
        gid_map: GID mapping in host:instance format (e.g. "1000:0").
        mount_type: Mount type: "classic" or "native".
    """
    args = ["mount", source, target]
    if uid_map:
        args.extend(["--uid-map", uid_map])
    if gid_map:
        args.extend(["--gid-map", gid_map])
    if mount_type:
        args.extend(["--type", mount_type])
    return await run_multipass(*args)


@mcp.tool()
async def umount_directory(mount_path: str) -> str:
    """Unmount a directory from a Multipass instance.

    Args:
        mount_path: Mount point to remove (instance:path format or host path).
    """
    return await run_multipass("umount", mount_path)


# ---------------------------------------------------------------------------
# Tools — Snapshots
# ---------------------------------------------------------------------------


@mcp.tool()
async def snapshot(
    instance: str,
    name: str = "",
    comment: str = "",
) -> str:
    """Create a snapshot of a stopped Multipass instance.

    The instance must be stopped before taking a snapshot.

    Args:
        instance: Instance name.
        name: Snapshot name. Auto-generated if omitted.
        comment: Description or comment for the snapshot.
    """
    args = ["snapshot", instance]
    if name:
        args.extend(["--name", name])
    if comment:
        args.extend(["--comment", comment])
    return await run_multipass(*args)


@mcp.tool()
async def restore(
    instance: str,
    snapshot: str,
    destructive: bool = False,
) -> str:
    """Restore a Multipass instance to a snapshot.

    The instance must be stopped. Use destructive=True to discard current state.

    Args:
        instance: Instance name.
        snapshot: Snapshot name to restore to.
        destructive: Discard current instance state (required if state changed since snapshot).
    """
    args = ["restore", f"{instance}.{snapshot}"]
    if destructive:
        args.append("--destructive")
    return await run_multipass(*args)


@mcp.tool()
async def clone(source_name: str, name: str = "") -> str:
    """Clone a stopped Multipass instance into an independent copy.

    The source instance must be stopped.

    Args:
        source_name: Name of the instance to clone.
        name: Name for the new cloned instance. Auto-generated if omitted.
    """
    args = ["clone", source_name]
    if name:
        args.extend(["--name", name])
    return await run_multipass(*args)


# ---------------------------------------------------------------------------
# Tools — Configuration
# ---------------------------------------------------------------------------


@mcp.tool()
async def get_config(key: str = "") -> str:
    """Get Multipass configuration settings.

    Args:
        key: Specific setting key (e.g. "local.driver"). Omit to get all settings.
    """
    args = ["get"]
    if key:
        args.append(key)
    else:
        args.append("--keys")
    return await run_multipass(*args)


@mcp.tool()
async def set_config(key: str, value: str) -> str:
    """Set a Multipass configuration value.

    Args:
        key: Setting key (e.g. "local.driver").
        value: New value to set.
    """
    return await run_multipass("set", f"{key}={value}")


# ---------------------------------------------------------------------------
# Tools — System
# ---------------------------------------------------------------------------


@mcp.tool()
async def purge() -> str:
    """Permanently delete all trashed (deleted) Multipass instances.

    This is irreversible. All instances previously deleted with 'delete' will be
    permanently removed.
    """
    return await run_multipass("purge")


@mcp.tool()
async def authenticate(passphrase: str) -> str:
    """Authenticate with the Multipass service using a passphrase.

    Args:
        passphrase: Authentication passphrase.
    """
    return await run_multipass("authenticate", passphrase)


# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------


def main():
    mcp.run(transport="stdio")


if __name__ == "__main__":
    main()
