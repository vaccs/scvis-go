# SCVis - Software Code Visualization System

SCVis is a web-based application for visualizing and analyzing C source code using the VACCS (Vulnerability Analysis of C Code Structure) analysis system. It provides an interactive interface for uploading C files, running security analysis, and visualizing the results.

## Features

- 🔍 **C Code Analysis**: Upload and analyze C source files for security vulnerabilities
- 📊 **Interactive Visualization**: View analysis results with detailed execution traces
- 👥 **User Management**: Multi-user support with administrator privileges
- 🔐 **Secure Authentication**: Bcrypt password hashing and session management
- 📁 **Project Management**: Organize analyses into projects
- 🚀 **Production Ready**: Systemd service integration with SSL/TLS support
- 📝 **Comprehensive Logging**: Structured logging with automatic rotation

## Quick Start

### 1. Database Setup

Set up the MySQL database (easiest method):

```bash
# Run the automated database installer
./install_database.sh
```

For manual setup or detailed instructions, see [DATABASE.md](DATABASE.md).

### 2. Build and Run

Development mode:
```bash
# Quick build and run
./quick-build.sh

# Or manually
go run ./cmd/web
```

Production deployment:
```bash
# Full build with Flutter frontend
./build.sh

# Install as systemd service (see BUILD.md)
```

### 3. Access the Application

- **Development**: https://localhost:4000
- **Production**: https://your-domain.com

Default admin credentials (if created during setup):
- Email: admin@localhost
- Password: admin123 (change immediately)

## Project Structure

```
scvis-go/
├── cmd/web/                 # Main application
│   ├── main.go             # Application entry point
│   ├── handlers.go         # HTTP handlers
│   ├── routes.go           # URL routing
│   └── ...                 # Other web components
├── internal/               # Internal packages
│   ├── models/            # Database models
│   ├── assert/            # Assertion helpers
│   ├── validator/         # Input validation
│   └── vaccs/             # VACCS integration
├── ui/                    # User interface
│   ├── html/             # HTML templates
│   └── static/           # Static assets (CSS, JS)
├── visual/               # Flutter web build output
├── scripts/              # Utility scripts
│   ├── create_user.go    # User creation tool
│   └── README.md         # Scripts documentation
├── log/                  # Application logs
├── tls/                  # SSL certificates (development)
├── input/                # Sample C files for testing
├── install_database.sh   # Database setup script
├── validate_database.sh  # Database validation script
├── build.sh             # Production build script
├── quick-build.sh       # Development build script
├── BUILD.md             # Build system documentation
├── DATABASE.md          # Database setup guide
└── SECURITY_MIGRATION.md # Security implementation notes
```

## Documentation

- **[Database Setup](DATABASE.md)** - Complete database installation and configuration guide
- **[Build System](BUILD.md)** - Building and deployment instructions
- **[Security Migration](SECURITY_MIGRATION.md)** - Security implementation details
- **[Scripts Documentation](scripts/README.md)** - User management and utility scripts

## Technology Stack

### Backend
- **Go 1.19+** - Primary application language
- **MySQL 8.0+** - Database backend
- **VACCS** - C code analysis engine

### Frontend
- **Flutter Web** - Modern web UI framework
- **Go Templates** - Server-side rendering
- **HTML/CSS/JavaScript** - Core web technologies

### Infrastructure
- **Systemd** - Service management (Linux)
- **Let's Encrypt** - SSL/TLS certificates
- **Nginx** - Reverse proxy (optional)

## Development Setup

### Prerequisites

```bash
# Go 1.19+
sudo apt install golang-go

# MySQL
sudo apt install mysql-server mysql-client

# Flutter (for frontend development)
# Follow instructions at https://flutter.dev/docs/get-started/install
```

### Environment Configuration

1. **Database Setup**:
   ```bash
   ./install_database.sh
   ```

2. **Environment Variables** (`.env`):
   ```bash
   DB_HOST=localhost
   DB_USER=scvis
   DB_PASSWD=your_password
   DB_NAME=scvis_local
   PORT=4000
   ```

3. **Development Certificates**:
   ```bash
   # Generate self-signed certificates for HTTPS
   openssl req -x509 -newkey rsa:4096 -keyout tls/localhost-key.pem -out tls/localhost.pem -days 365 -nodes
   ```

### Running in Development

```bash
# Start the application
go run ./cmd/web

# Or use the quick build script
./quick-build.sh
```

Visit https://localhost:4000 to access the application.

## Production Deployment

### 1. Prepare the Server

```bash
# Install dependencies
sudo apt update
sudo apt install mysql-server golang-go

# Create application user
sudo useradd -r -s /bin/false -d /opt/scvis www-data
```

### 2. Build and Deploy

```bash
# Clone and build
git clone <repository-url> /opt/scvis
cd /opt/scvis

# Set up database
./install_database.sh

# Build application
./build.sh

# Configure production environment
sudo mkdir -p /etc/scvisual
sudo cp scvisual.conf.example /etc/scvisual/scvisual.conf
sudo chown root:www-data /etc/scvisual/scvisual.conf
sudo chmod 640 /etc/scvisual/scvisual.conf
```

### 3. Install as System Service

```bash
# Copy and enable service
sudo cp scvisual.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable scvisual
sudo systemctl start scvisual
```

### 4. SSL/TLS Setup

For production with Let's Encrypt:

```bash
# Install certbot
sudo apt install certbot

# Obtain certificate
sudo certbot certonly --standalone -d your-domain.com

# Update service configuration to use real certificates
```

## User Management

### Creating Users

```bash
# Interactive user creation
cd scripts
go run create_user.go
```

### Managing Users via Database

```sql
-- Create admin user
INSERT INTO users (name, email, hashed_password, administrator) 
VALUES ('Admin', 'admin@example.com', '$2b$12$hash_here', true);

-- List all users
SELECT id, name, email, administrator FROM users;

-- Update user permissions
UPDATE users SET administrator = true WHERE email = 'user@example.com';
```

## API Endpoints

### Authentication
- `GET /user/signup` - User registration form
- `POST /user/signup` - Process user registration
- `GET /user/login` - Login form
- `POST /user/login` - Process login
- `POST /user/logout` - Logout

### Projects
- `GET /` - Project dashboard
- `GET /project/create` - Create project form
- `POST /project/create` - Process project creation
- `GET /project/view/{id}` - View project details
- `POST /project/analyze/{id}` - Run analysis

### File Management
- `POST /cfile/upload` - Upload C files
- `GET /cfile/view/{id}` - View file content

## Security Features

- **Bcrypt Password Hashing** - Cost factor 12
- **CSRF Protection** - Secure session management
- **SQL Injection Prevention** - Prepared statements
- **TLS/SSL Encryption** - HTTPS only
- **Session Security** - Secure cookie flags
- **Input Validation** - Comprehensive input sanitization
- **Privilege Separation** - Runs as www-data user
- **Security Headers** - HSTS, CSP, etc.

## Monitoring and Logging

### Log Files
- **Application Logs**: `/var/log/scvisual/scvisual.log`
- **Error Logs**: Included in main log with ERROR level
- **System Logs**: `journalctl -u scvisual`

### Log Rotation
- Automatic rotation at 10MB (configurable)
- Maintains 5 backup files
- Atomic rotation with zero downtime

### Monitoring Commands

```bash
# Service status
sudo systemctl status scvisual

# Live logs
sudo tail -f /var/log/scvisual/scvisual.log

# Service restart
sudo systemctl restart scvisual

# Check for errors
sudo journalctl -u scvisual --since "1 hour ago"
```

## Troubleshooting

### Common Issues

**Database Connection Failed**:
```bash
# Check MySQL service
sudo systemctl status mysql

# Validate database setup
./validate_database.sh

# Test connection manually
mysql -u scvis -p scvis_local
```

**Permission Denied**:
```bash
# Check file ownership
sudo chown -R www-data:www-data /opt/scvis

# Check service permissions
sudo systemctl status scvisual
```

**SSL Certificate Issues**:
```bash
# Check certificate files
sudo ls -la /etc/letsencrypt/live/your-domain.com/

# Renew certificates
sudo certbot renew
```

### Debug Mode

Enable debug logging by setting environment variable:
```bash
export LOG_LEVEL=debug
go run ./cmd/web
```

### Health Checks

```bash
# Database validation
./validate_database.sh

# Application health
curl -k https://localhost:4000/

# Service health
sudo systemctl is-active scvisual
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

### Development Guidelines

- Follow Go coding standards
- Add documentation for new features
- Update relevant README files
- Test database migrations
- Ensure security best practices

## License

[Your License Here]

## Support

For support and questions:
- Check the documentation in this repository
- Review log files for error details
- Ensure database is properly configured
- Verify all dependencies are installed

## Version History

### Latest
- ✅ Complete database automation with install_database.sh
- ✅ Database validation with validate_database.sh
- ✅ Comprehensive build system (build.sh, quick-build.sh)
- ✅ Production deployment with systemd integration
- ✅ Security hardening and SSL/TLS support
- ✅ Log rotation and monitoring capabilities
- ✅ User management tools and documentation

### Previous
- Department column removal and code cleanup
- Flutter web frontend integration
- VACCS analysis engine integration
- Multi-user authentication system
- Project and file management
