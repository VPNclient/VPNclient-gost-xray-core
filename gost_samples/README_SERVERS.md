# GLess and GMess Server Setup

This document provides instructions for running GLess and GMess servers with the provided configuration files.

## Prerequisites

1. Build the xray binary:
   ```bash
   ./build.sh
   ```

2. Generate SSL certificates (for TLS):
   ```bash
   # Generate self-signed certificate for testing
   openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
   ```

## Server Configurations

### GLess Server

**Configuration File:** `gless_server_config.json`
**Port:** 443 (HTTPS)
**Protocol:** GLess with GOST cryptography

**Features:**
- GOST-based encryption (gost-rprx-vision flow)
- TLS encryption
- Fallback support
- User authentication with UUID

**Run Command:**
```bash
./run_gless_server.sh
```

### GMess Server

**Configuration File:** `gmess_server_config.json`
**Port:** 444 (HTTPS)
**Protocol:** GMess with GOST cryptography

**Features:**
- GOST-based encryption (gost-28147 security)
- TLS encryption
- User authentication with UUID

**Run Command:**
```bash
./run_gmess_server.sh
```

## Client Configurations

### GLess Client

**Configuration File:** `gless_client_config.json`
**Local Ports:** 
- SOCKS: 1080
- HTTP: 1081

**Features:**
- SOCKS5 proxy support
- HTTP proxy support
- GOST encryption

### GMess Client

**Configuration File:** `gmess_client_config.json`
**Local Ports:**
- SOCKS: 1080
- HTTP: 1081

**Features:**
- SOCKS5 proxy support
- HTTP proxy support
- GOST encryption

## Configuration Customization

### Server Configuration

1. **Change Ports:** Modify the `port` field in the server configs
2. **Update UUIDs:** Generate new UUIDs for user authentication
3. **SSL Certificates:** Replace `cert.pem` and `key.pem` with your certificates
4. **Server IP:** Update `your-server-ip` in client configs

### Client Configuration

1. **Server Address:** Change `your-server-ip` to your actual server IP
2. **Domain:** Update `your-domain.com` to match your SSL certificate
3. **UUIDs:** Ensure client UUIDs match server UUIDs

## Security Notes

1. **UUID Generation:** Use secure UUID generators for production
2. **SSL Certificates:** Use valid SSL certificates for production
3. **Firewall:** Configure firewall to allow server ports
4. **Logs:** Monitor access.log and error.log for issues

## Troubleshooting

1. **Port Already in Use:** Change port numbers in configuration
2. **SSL Errors:** Ensure certificate files exist and are valid
3. **Connection Refused:** Check firewall settings and server status
4. **Authentication Failed:** Verify UUIDs match between client and server

## Example Usage

### Start GLess Server
```bash
./run_gless_server.sh
```

### Start GMess Server
```bash
./run_gmess_server.sh
```

### Test Connection (Client)
```bash
# Test SOCKS proxy
curl --socks5 127.0.0.1:1080 http://httpbin.org/ip

# Test HTTP proxy
curl -x http://127.0.0.1:1081 http://httpbin.org/ip
```

## File Structure

```
├── run_gless_server.sh          # GLess server runner
├── run_gmess_server.sh          # GMess server runner
├── gless_server_config.json     # GLess server configuration
├── gmess_server_config.json     # GMess server configuration
├── gless_client_config.json     # GLess client configuration
├── gmess_client_config.json     # GMess client configuration
├── cert.pem                     # SSL certificate (generate)
├── key.pem                      # SSL private key (generate)
└── xray                         # Compiled binary
``` 