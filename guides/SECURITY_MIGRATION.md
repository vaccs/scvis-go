# SCVISUAL Security Migration: Running as www-data

This document describes the changes made to run SCVISUAL as the `www-data` user instead of root for improved security.

## Changes Made

### 1. Code Changes
- **SSL Certificate Paths**: Uses the existing Let's Encrypt certificate paths accessible by the `ssl-cert` group:
  - Certificate: `/etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem`
  - Private Key: `/etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem`

### 2. Service Configuration
- **User/Group**: Service now runs as `www-data:www-data`
- **Security Hardening**: Added systemd security features:
  - `NoNewPrivileges=true`
  - `PrivateTmp=true`
  - `ProtectSystem=strict`
  - `ProtectHome=true`
  - Limited capabilities to only `CAP_NET_BIND_SERVICE`

### 3. File Structure
```
/opt/scvisual/          # Application directory (owned by www-data)
├── scvisual            # Binary
├── visual/             # Static files
└── ui/                 # UI templates

/etc/scvisual/          # Configuration directory (root:www-data, 750)
└── scvisual.conf       # Configuration file (root:www-data, 640)

/var/log/scvisual/      # Log directory (www-data:www-data, 755)

/etc/letsencrypt/live/vaccs.ceas.wmich.edu/  # SSL certificates (accessible by ssl-cert group)
├── fullchain.pem       # Certificate chain
└── privkey.pem         # Private key
```

## Migration Steps

### 1. Stop Current Service
```bash
sudo systemctl stop scvisual
```

### 2. Run Deployment Script
```bash
sudo ./deploy.sh
```

This script will:
- Create necessary directories with proper ownership
- Build and install the application
- Configure SSL certificate access
- Install the systemd service file

### 3. Update SSL Certificates
The service will use the existing Let's Encrypt certificates that are already accessible by the `ssl-cert` group. No copying is necessary - just ensure `www-data` is in the `ssl-cert` group:
```bash
sudo usermod -a -G ssl-cert www-data
```

### 4. Start New Service
```bash
sudo systemctl enable scvisual
sudo systemctl start scvisual
```

### 5. Verify Operation
```bash
sudo systemctl status scvisual
sudo journalctl -u scvisual -f
```

## Security Benefits

1. **Principle of Least Privilege**: Service no longer runs as root
2. **Limited Capabilities**: Only has network binding capability
3. **Filesystem Protection**: Limited write access to specific directories
4. **Process Isolation**: Cannot access user home directories or create new privileges
5. **SSL Certificate Access**: Controlled through group membership rather than root access

## Troubleshooting

### SSL Certificate Issues
- Ensure `www-data` is in the `ssl-cert` group: `groups www-data`
- Check certificate file permissions: `ls -la /etc/letsencrypt/live/vaccs.ceas.wmich.edu/`
- Verify group membership: `sudo usermod -a -G ssl-cert www-data`
- Test certificate access: `sudo -u www-data cat /etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem > /dev/null`

### Permission Issues
- Check application directory ownership: `ls -la /opt/scvisual`
- Check log directory access: `ls -la /var/log/scvisual`
- Test writing to log directory: `sudo -u www-data touch /var/log/scvisual/test.log`

### Service Issues
- Check service status: `sudo systemctl status scvisual`
- View detailed logs: `sudo journalctl -u scvisual -n 50`
- Check configuration: `sudo -u www-data /opt/scvisual/scvisual --help`
