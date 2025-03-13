#!/bin/bash

# Define port range
MIN_PORT=10000
MAX_PORT=10010

echo "Closing all open ports between $MIN_PORT and $MAX_PORT..."

# Loop through each port in range
for ((port=MIN_PORT; port<=MAX_PORT; port++)); do
    # Find process using the port
    pid=$(lsof -ti :$port)  # Works on macOS and Linux with lsof

    # Alternative method using netstat or ss
    # pid=$(netstat -tulpn 2>/dev/null | grep ":$port " | awk '{print $7}' | cut -d'/' -f1)
    # pid=$(ss -tulpn 2>/dev/null | grep ":$port " | awk '{print $7}' | cut -d',' -f2 | cut -d'=' -f2)

    if [[ -n "$pid" ]]; then
        echo "Closing port $port (PID: $pid)"
        kill -9 $pid
    else
        echo "Port $port is not in use."
    fi
done

echo "All specified ports have been checked."
