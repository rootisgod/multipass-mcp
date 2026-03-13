#!/bin/bash
set -euo pipefail

# Health check script for Multipass VMs
# Transfer into a VM and execute to verify basic health.
# Returns a JSON summary.

apt_status="ok"
dns_status="ok"
disk_free_gb="0"
mem_free_mb="0"

# Check apt-get update
if ! sudo apt-get update -qq > /dev/null 2>&1; then
    apt_status="fail"
fi

# Check DNS resolution
if ! ping -c1 -W5 archive.ubuntu.com > /dev/null 2>&1; then
    dns_status="fail"
fi

# Check free disk space (root filesystem)
disk_free_kb=$(df / | awk 'NR==2 {print $4}')
disk_free_gb=$(echo "scale=1; $disk_free_kb / 1048576" | bc)

# Check free memory
mem_free_kb=$(grep MemAvailable /proc/meminfo | awk '{print $2}')
mem_free_mb=$(echo "scale=0; $mem_free_kb / 1024" | bc)

# Output JSON summary
cat <<EOF
{"apt": "$apt_status", "dns": "$dns_status", "disk_free_gb": $disk_free_gb, "mem_free_mb": $mem_free_mb}
EOF
