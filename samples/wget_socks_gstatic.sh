#!/bin/bash

# Download gstatic file via local SOCKS5 proxy (127.0.0.1:1080)
URL="https://www.gstatic.com/images/branding/product/1x/googleg_32dp.png"
OUT="googleg_32dp.png"

# Use curl instead of wget for better SOCKS5 support
curl --socks5 127.0.0.1:1080 --socks5-hostname 127.0.0.1:1080 -k -o "$OUT" "$URL"

echo "Download completed: $OUT"
