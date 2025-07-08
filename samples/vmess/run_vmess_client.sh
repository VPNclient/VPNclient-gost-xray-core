#!/bin/bash

# VMess Client Runner Script
# This script runs the VMess client with SOCKS5 proxy

echo "Starting VMess Client..."

# Check if xray binary exists
if [ ! -f "./xray" ]; then
    echo "Error: xray binary not found. Please copy it from the main directory."
    exit 1
fi

# Check if config file exists
if [ ! -f "./vmess_client_config.json" ]; then
    echo "Error: vmess_client_config.json not found."
    exit 1
fi

# Run the VMess client
echo "VMess client starting..."
echo "SOCKS5 proxy will be available on 127.0.0.1:1080"
echo "Press Ctrl+C to stop the client"

./xray -config=vmess_client_config.json 