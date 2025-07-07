#!/bin/bash

# Simple SOCKS5 test script
echo "Testing SOCKS5 proxy connection..."

# Test with a simple HTTP request
curl --socks5 127.0.0.1:1080 --socks5-hostname 127.0.0.1:1080 -k -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://httpbin.org/ip

echo "Test completed." 