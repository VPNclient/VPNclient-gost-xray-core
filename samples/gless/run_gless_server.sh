#!/bin/bash

# GLess Server Runner Script
# This script runs a GLess server with the provided configuration

CONFIG_FILE="gless_server_config.json"

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

echo "Starting GLess server with configuration: $CONFIG_FILE"
echo "Press Ctrl+C to stop the server"

# Run the server
./xray run -config "$CONFIG_FILE" 