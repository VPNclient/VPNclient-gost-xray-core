#!/bin/bash

# SSL Certificate Generator for GLess/GMess Servers
# This script generates self-signed SSL certificates for testing

echo "Generating SSL certificates for GLess/GMess servers..."

# Check if OpenSSL is available
if ! command -v openssl &> /dev/null; then
    echo "Error: OpenSSL is not installed!"
    echo "Please install OpenSSL first."
    exit 1
fi

# Generate private key and certificate
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

if [ $? -eq 0 ]; then
    echo "SSL certificates generated successfully!"
    echo "Files created:"
    echo "  - cert.pem (certificate)"
    echo "  - key.pem (private key)"
    echo ""
    echo "Note: These are self-signed certificates for testing only."
    echo "For production use, obtain valid SSL certificates from a CA."
else
    echo "Error: Failed to generate SSL certificates!"
    exit 1
fi 