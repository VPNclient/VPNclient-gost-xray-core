#!/bin/bash

# VMess Server Runner Script
# This script runs the VMess server with GOST encryption

echo "Starting VMess Server..."

# Check if xray binary exists
if [ ! -f "./xray" ]; then
    echo "Error: xray binary not found. Please copy it from the main directory."
    exit 1
fi

# Check if config file exists
if [ ! -f "./vmess_server_config.json" ]; then
    echo "Error: vmess_server_config.json not found."
    exit 1
fi

# Check if SSL certificates exist
if [ ! -f "./server.crt" ] || [ ! -f "./server.key" ]; then
    echo "Warning: SSL certificates not found. Generating GOST certificates..."
    ./gost_generate_cert.sh
fi

# Run the VMess server
echo "VMess server starting on port 8443..."
echo "Client ID: b831381d-6324-4d53-ad4f-8cda48b30811"
echo "Press Ctrl+C to stop the server"

./xray -config=vmess_server_config.json 