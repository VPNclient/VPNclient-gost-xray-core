#!/bin/bash

# GLess Client Runner Script
# This script runs a GLess client with the provided configuration

CONFIG_FILE="gless_client_config.json"

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Configuration file $CONFIG_FILE not found!"
    echo "Please create the configuration file first."
    exit 1
fi

# Check if xray binary exists
if [ ! -f "./xray" ]; then
    echo "Error: xray binary not found!"
    echo "Please run ./build.sh first to build the binary."
    exit 1
fi

echo "Starting GLess client with configuration: $CONFIG_FILE"
echo "Press Ctrl+C to stop the client"
echo ""
echo "Client will be available on:"
echo "  - SOCKS5: 127.0.0.1:1080"
echo "  - HTTP: 127.0.0.1:1081"
echo ""

# Run the client
./xray run -config "$CONFIG_FILE" 