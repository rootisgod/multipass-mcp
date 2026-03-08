# CLAUDE.md

## Project Overview

Multipass MCP Server — exposes the Multipass CLI as an MCP server with 18 tools and 6 resources. Lets AI assistants manage Ubuntu VMs.

## Tech Stack

- Python 3.10+, single-module server
- FastMCP from `mcp[cli]` package
- Hatchling build system
- Async subprocess calls to `multipass` CLI

## Project Structure

```
src/multipass_mcp/
├── __init__.py     # empty package marker
└── server.py       # entire server: helpers, 6 resources, 18 tools, main()
pyproject.toml      # build config, deps, entry point
```

All logic lives in `server.py`. No need for multiple modules at this scale.

## Key Patterns

- `run_multipass(*args)` — async subprocess helper, raises `RuntimeError` on failure
- `run_multipass_json(*args)` — same but appends `--format json` and parses result
- Resources use `@mcp.resource()`, tools use `@mcp.tool()`
- `name="all"` on lifecycle tools translates to `--all` CLI flag
- `exec_command` takes `command: list[str]` (not a shell string) to prevent injection
- Timeouts: 300s default, 600s for `launch`

## Development Commands

```bash
# Install in dev mode
uv venv && uv pip install -e .

# Test with MCP Inspector
mcp dev src/multipass_mcp/server.py

# Register with Claude Code
claude mcp add multipass-mcp -- /path/to/.venv/bin/multipass-mcp
```

## Entry Point

`multipass-mcp` CLI → `multipass_mcp.server:main()` → `mcp.run(transport="stdio")`

## Verified Working

End-to-end tested: launch VM, install packages via exec_command, deploy files, delete+purge. All 18 tools and 6 resources confirmed registered (5 static + 1 template).

## Notes

- `exec_command` can run complex shell commands via `["bash", "-c", "..."]`
- CLI spinner/escape codes appear in tool output but don't affect functionality
- Empty string result from tools (e.g. delete, mount) means success
- Ubuntu 24.04 (Noble) is the current default LTS image on arm64
