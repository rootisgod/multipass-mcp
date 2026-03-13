#!/bin/bash
set -euo pipefail

# Set up SSH authorized keys for the current user.
# Usage: bash setup_ssh_keys.sh "ssh-rsa AAAA... user@host"

if [ $# -lt 1 ]; then
    echo "Usage: $0 <public_key>"
    exit 1
fi

PUBLIC_KEY="$1"

# Create .ssh directory if it doesn't exist
mkdir -p ~/.ssh
chmod 700 ~/.ssh

# Append the key if not already present
if ! grep -qF "$PUBLIC_KEY" ~/.ssh/authorized_keys 2>/dev/null; then
    echo "$PUBLIC_KEY" >> ~/.ssh/authorized_keys
    echo "Key added to ~/.ssh/authorized_keys"
else
    echo "Key already present in ~/.ssh/authorized_keys"
fi

# Set correct permissions
chmod 600 ~/.ssh/authorized_keys

echo "SSH key setup complete. Permissions verified."
