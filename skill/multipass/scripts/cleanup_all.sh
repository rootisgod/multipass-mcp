#!/bin/bash
set -euo pipefail

# HOST-SIDE cleanup script — stops, deletes, and purges ALL Multipass instances.
# Run this on the host machine, NOT inside a VM.
# WARNING: This is destructive and irreversible!

echo "Counting running instances..."
count=$(multipass list --format json 2>/dev/null | python3 -c "
import sys, json
data = json.load(sys.stdin)
instances = data.get('list', [])
print(len(instances))
" 2>/dev/null || echo "0")

if [ "$count" = "0" ]; then
    echo "No instances found. Nothing to clean up."
    exit 0
fi

echo "Found $count instance(s). Stopping all..."
multipass stop --all 2>/dev/null || true

echo "Deleting all instances..."
multipass delete --all 2>/dev/null || true

echo "Purging deleted instances..."
multipass purge 2>/dev/null || true

echo "Cleanup complete. Removed $count instance(s)."
