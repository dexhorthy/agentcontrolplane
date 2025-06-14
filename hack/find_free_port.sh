#!/bin/bash
# find_free_port.sh - Find an available port in the specified range
# Usage: ./find_free_port.sh [start_port] [end_port]

set -euo pipefail

# Default port range
START_PORT=${1:-10000}
END_PORT=${2:-11000}

# Function to check if a port is available
is_port_available() {
    local port=$1
    # Check if port is in use using netstat or ss
    if command -v ss >/dev/null 2>&1; then
        ! ss -tuln | grep -q ":${port} "
    elif command -v netstat >/dev/null 2>&1; then
        ! netstat -tuln | grep -q ":${port} "
    else
        # Fallback: try to bind to the port
        ! nc -z localhost "$port" 2>/dev/null
    fi
}

# Find first available port in range
for port in $(seq "$START_PORT" "$END_PORT"); do
    if is_port_available "$port"; then
        echo "$port"
        exit 0
    fi
done

# No available port found
echo "Error: No available port found in range $START_PORT-$END_PORT" >&2
exit 1