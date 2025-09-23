#!/bin/bash
# Entrypoint script to ensure proper ownership of mounted volumes
# This script runs as root at container startup to fix permissions before switching to testuser

set -e

# Fix ownership of the cache directory if it exists and is owned by root
if [[ -d "/home/testuser/.cache" ]]; then
    # Check if the directory is owned by root (which happens with Docker volume mounts)
    if [[ "$(stat -c %U /home/testuser/.cache)" == "root" ]]; then
        echo "🔧 Fixing ownership of /home/testuser/.cache directory..."
        chown -R testuser:testuser /home/testuser/.cache
    fi
fi

# Ensure the entire home directory has correct ownership
if [[ "$(stat -c %U /home/testuser)" == "root" ]]; then
    echo "🔧 Fixing ownership of /home/testuser directory..."
    chown testuser:testuser /home/testuser
fi

# If no command is provided, start bash as testuser
if [[ $# -eq 0 ]]; then
    exec su - testuser
else
    # Execute the provided command as testuser
    exec su - testuser -c "$*"
fi
