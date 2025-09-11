# Configuration Updates - September 11, 2025

## Summary

Updated the SCVis application to use standardized `.env` configuration files for both development and production deployments, replacing the previous mixed approach of `.env` for development and `/etc/scvisual/scvisual.conf` for production.

## Changes Made

### 1. Unified Configuration System

**Previous Approach:**
- Development: `.env` file in project root
- Production: `/etc/scvisual/scvisual.conf` file in system directory

**New Approach:**
- Development: `.env` file in project root
- Production: `.env` file in `/opt/scvis/` directory

### 2. Environment Variables Added

**New Configuration:**
- `LOCAL_SERVER` - Explicit server mode setting
  - `true` - Development/local server mode (uses development defaults)
  - `false` - Production server mode (uses production defaults)
  - Auto-detected if not explicitly set

**Existing Configuration (Enhanced):**
- `LOG_DIR` - Configurable log directory path
- `LOG_MAX_SIZE_MB` - Configurable log rotation size threshold
- `LOG_MAX_BACKUPS` - Configurable number of backup files to keep
- `SSL_CERT_PATH` - Configurable SSL certificate file path
- `SSL_KEY_PATH` - Configurable SSL private key file path

### 3. Files Modified

#### Application Code
- **`cmd/web/main.go`** - Updated server mode detection and configuration loading
  - Replaced hardcoded `/etc/scvisual/scvisual.conf` detection
  - Added `LOCAL_SERVER` environment variable support
  - Implemented auto-detection fallback based on directory structure
  - Changed production config path from `/etc/scvisual/scvisual.conf` to `/opt/scvis/.env`

#### Configuration Files
- **`.env`** - Added `LOCAL_SERVER=true` setting for development
- **`.env.example`** - Added `LOCAL_SERVER` configuration option with documentation
- **`production.env`** - New production configuration template (replaces `scvisual.conf.example`)

#### Installation and Documentation
- **`install_database.sh`** - Updated to generate new configuration format
  - Creates `production.env` instead of `scvisual.conf.example`
  - Updated deployment instructions to use `/opt/scvis/.env`
  - Added `LOCAL_SERVER` setting to generated files
- **`README.md`** - Updated configuration documentation
  - Changed production config path references
  - Added `LOCAL_SERVER` configuration option documentation
  - Updated deployment instructions
- **`DATABASE.md`** - Updated production configuration examples

### 4. Behavior Changes

#### Configuration Loading
- **Before:** Checked for `/etc/scvisual/scvisual.conf` existence to determine production mode
- **After:** 
  1. Checks `LOCAL_SERVER` environment variable for explicit mode setting
  2. Falls back to directory structure detection
  3. Uses `/opt/scvis/.env` for production configuration

#### Server Mode Detection
- **Before:** Production mode triggered by `/etc/scvisual/scvisual.conf` file existence
- **After:** 
  - Explicit: `LOCAL_SERVER=true/false` environment variable
  - Auto-detection: Based on presence of `./scvis-go` (dev) vs `/opt/scvis/.env` (prod)
  - Default: Local server mode for backward compatibility

### 5. Migration Guide

#### For Development
No changes required - existing `.env` files continue to work. Optionally add `LOCAL_SERVER=true` for explicit mode setting.

#### For Production Deployments
Replace the old system configuration:

**Old Approach:**
```bash
sudo mkdir -p /etc/scvisual
sudo cp scvisual.conf.example /etc/scvisual/scvisual.conf
sudo chown root:www-data /etc/scvisual/scvisual.conf
sudo chmod 640 /etc/scvisual/scvisual.conf
```

**New Approach:**
```bash
sudo mkdir -p /opt/scvis
sudo cp production.env /opt/scvis/.env
sudo chown root:www-data /opt/scvis/.env
sudo chmod 640 /opt/scvis/.env
```

### 6. Benefits

1. **Consistency** - Both development and production use `.env` format
2. **Standardization** - Follows common `.env` configuration patterns
3. **Flexibility** - Explicit mode control via `LOCAL_SERVER` variable
4. **Maintainability** - Single configuration file format to manage
5. **Portability** - Configuration is self-contained in application directory

### 7. Backward Compatibility

✅ **Development environments** - Fully backward compatible
⚠️  **Production deployments** - Requires migration from `/etc/scvisual/scvisual.conf` to `/opt/scvis/.env`

### 8. Default Values

| Variable | Development Default | Production Default |
|----------|--------------------|--------------------|
| `LOCAL_SERVER` | `true` | `false` |
| `LOG_DIR` | `./log/scvisual` | `/var/log/scvisual` |
| `LOG_MAX_SIZE_MB` | `1` | `10` |
| `LOG_MAX_BACKUPS` | `5` | `5` |
| `SSL_CERT_PATH` | `./tls/localhost.pem` | `/etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem` |
| `SSL_KEY_PATH` | `./tls/localhost-key.pem` | `/etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem` |

### 9. Configuration File Paths

| Environment | Configuration File | 
|-------------|-------------------|
| Development | `./env` |
| Production | `/opt/scvis/.env` |

### 10. Validation

✅ Application builds successfully with new configuration logic
✅ Local development mode works with existing `.env` files
✅ Production mode correctly attempts to load from `/opt/scvis/.env`
✅ Explicit mode setting via `LOCAL_SERVER` environment variable works
✅ Auto-detection fallback logic functions properly
✅ Installation script generates correct configuration templates

## Usage Examples

### Explicit Mode Setting
```bash
# Force local server mode
LOCAL_SERVER=true ./scvis

# Force production server mode  
LOCAL_SERVER=false ./scvis
```

### Development Configuration (`.env`)
```bash
DB_HOST=localhost
DB_USER=scvis
DB_PASSWD=dev_password
DB_NAME=scvis_local
PORT=4000
LOCAL_SERVER=true

# Log Configuration
LOG_DIR=./log/scvisual
LOG_MAX_SIZE_MB=1
LOG_MAX_BACKUPS=5

# SSL Certificate Configuration
SSL_CERT_PATH=./tls/localhost.pem
SSL_KEY_PATH=./tls/localhost-key.pem
```

### Production Configuration (`/opt/scvis/.env`)
```bash
DB_HOST=localhost
DB_USER=scvis
DB_PASSWD=production_password
DB_NAME=scvis
PORT=4000
LOCAL_SERVER=false

# Log Configuration
LOG_DIR=/var/log/scvisual
LOG_MAX_SIZE_MB=10
LOG_MAX_BACKUPS=5

# SSL Certificate Configuration
SSL_CERT_PATH=/etc/letsencrypt/live/your-domain.com/fullchain.pem
SSL_KEY_PATH=/etc/letsencrypt/live/your-domain.com/privkey.pem
```
