# Quick Start Guide

## 1. Build the Project
```bash
./build.sh
```

## 2. Generate SSL Certificates
```bash
./generate_cert.sh
```

## 3. Run Servers

### GLess Server
```bash
./run_gless_server.sh
```

### GMess Server
```bash
./run_gmess_server.sh
```

## 4. Test with Client

### GLess Client
```bash
./xray run -config gless_client_config.json
```

### GMess Client
```bash
./xray run -config gmess_client_config.json
```

## 5. Test Connection
```bash
# Test SOCKS proxy
curl --socks5 127.0.0.1:1080 http://httpbin.org/ip

# Test HTTP proxy
curl -x http://127.0.0.1:1081 http://httpbin.org/ip
```

## Configuration Files
- `gless_server_config.json` - GLess server configuration
- `gmess_server_config.json` - GMess server configuration
- `gless_client_config.json` - GLess client configuration
- `gmess_client_config.json` - GMess client configuration

## Important Notes
1. Update `your-server-ip` in client configs with actual server IP
2. Update `your-domain.com` in client configs to match SSL certificate
3. Generate new UUIDs for production use
4. Use valid SSL certificates for production 