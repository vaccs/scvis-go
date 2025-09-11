# SCVis Database Installation Guide

This guide covers the complete database setup process for the SCVis application, including both development and production databases.

## Overview

The SCVis application uses MySQL as its database backend with the following structure:
- **Production Database**: `scvis` - Used for production deployment
- **Development Database**: `scvis_local` - Used for local development and testing
- **Database User**: Configurable user with access to both databases
- **Tables**: `users`, `projects`, `cfiles`, `sessions`

## Prerequisites

### System Requirements
- MySQL Server 8.0+ (or MariaDB 10.3+)
- MySQL client tools
- Python3 with bcrypt (optional, for password hashing)
- Bash shell

### Installation Commands

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install mysql-server mysql-client-core-8.0 python3 python3-bcrypt
sudo systemctl start mysql
sudo systemctl enable mysql
```

**RHEL/CentOS/Fedora:**
```bash
sudo yum install mysql-server mysql python3 python3-bcrypt
# or for newer versions:
sudo dnf install mysql-server mysql python3 python3-bcrypt
sudo systemctl start mysqld
sudo systemctl enable mysqld
```

### MySQL Initial Setup

If this is a fresh MySQL installation, secure it first:
```bash
sudo mysql_secure_installation
```

## Automated Installation

The easiest way to set up the database is using the automated installation script:

```bash
# Navigate to the project directory
cd /path/to/scvis-go

# Run the installation script
./install_database.sh
```

### Installation Process

The script will guide you through:

1. **MySQL Connection Check**: Verifies MySQL is running and accessible
2. **Credentials Input**: Prompts for MySQL root credentials
3. **Database Configuration**: Allows customization of database settings
4. **Database Creation**: Creates both production and development databases
5. **User Setup**: Creates a dedicated database user with appropriate permissions
6. **Table Creation**: Creates all required tables with proper relationships
7. **Sample Data**: Optionally creates an admin user for testing
8. **Configuration Files**: Generates environment configuration files

### Default Configuration

If you accept all defaults, the script creates:
- **Production DB**: `scvis`
- **Development DB**: `scvis_local`
- **Database User**: `scvis`
- **Host**: `localhost`
- **Port**: `4000`

## Manual Installation

If you prefer to set up the database manually, follow these steps:

### 1. Connect to MySQL

```bash
mysql -u root -p
```

### 2. Create Databases

```sql
CREATE DATABASE scvis CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
CREATE DATABASE scvis_local CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

### 3. Create Database User

```sql
CREATE USER 'scvis'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT ALL PRIVILEGES ON scvis.* TO 'scvis'@'localhost';
GRANT ALL PRIVILEGES ON scvis_local.* TO 'scvis'@'localhost';
FLUSH PRIVILEGES;
```

### 4. Create Tables

Use the development database for table creation:

```sql
USE scvis_local;

-- Users table
CREATE TABLE users (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  hashed_password char(60) NOT NULL,
  administrator boolean NOT NULL DEFAULT false,
  created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY users_uc_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Projects table
CREATE TABLE projects (
  id int NOT NULL AUTO_INCREMENT,
  title varchar(255) NOT NULL,
  userid int NOT NULL,
  cfileids json NOT NULL,
  analysis blob,
  created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  parameters varchar(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY fk_projects_userid (userid),
  CONSTRAINT fk_projects_userid FOREIGN KEY (userid) REFERENCES users (id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- C Files table
CREATE TABLE cfiles (
  id int NOT NULL AUTO_INCREMENT,
  filename varchar(255) NOT NULL,
  content text NOT NULL,
  projectid int NOT NULL,
  created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY fk_cfiles_projectid (projectid),
  CONSTRAINT fk_cfiles_projectid FOREIGN KEY (projectid) REFERENCES projects (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Sessions table
CREATE TABLE sessions (
  token char(43) NOT NULL,
  data blob NOT NULL,
  expiry timestamp(6) NOT NULL,
  PRIMARY KEY (token),
  KEY sessions_expiry_idx (expiry)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

### 5. Repeat for Production

```sql
USE scvis;
-- Run the same CREATE TABLE statements as above
```

## Configuration Files

### Development Configuration (`.env`)

Create a `.env` file in the project root:

```bash
DB_HOST=localhost
DB_USER=scvis
DB_PASSWD=your_secure_password
DB_NAME=scvis_local
PORT=4000

# Log Configuration
LOG_DIR=./log/scvisual
LOG_MAX_SIZE_MB=1
LOG_MAX_BACKUPS=5

# SSL Certificate Configuration
SSL_CERT_PATH=./tls/localhost.pem
SSL_KEY_PATH=./tls/localhost-key.pem
```

### Production Configuration

For production deployment, create `/opt/scvis/.env`:

```bash
# Create directory
sudo mkdir -p /opt/scvis

# Create configuration file
sudo tee /opt/scvis/.env << EOF
DB_HOST=localhost
DB_USER=scvis
DB_PASSWD=your_secure_password
DB_NAME=scvis
PORT=4000

# Server Mode Configuration
LOCAL_SERVER=false

# Log Configuration
LOG_DIR=/var/log/scvisual
LOG_MAX_SIZE_MB=10
LOG_MAX_BACKUPS=5

# SSL Certificate Configuration
SSL_CERT_PATH=/etc/letsencrypt/live/your-domain.com/fullchain.pem
SSL_KEY_PATH=/etc/letsencrypt/live/your-domain.com/privkey.pem
EOF

# Set proper permissions
sudo chown root:www-data /opt/scvis/.env
sudo chmod 640 /opt/scvis/.env
```

## User Management

### Creating Users

Use the provided user creation script:

```bash
cd scripts
go run create_user.go
```

Or create users manually with bcrypt-hashed passwords:

```sql
-- Example with bcrypt hash (cost 12)
INSERT INTO users (name, email, hashed_password, administrator) 
VALUES ('Admin User', 'admin@example.com', '$2b$12$hash_here', true);
```

### Password Requirements

- Minimum 8 characters
- Bcrypt hashing with cost 12
- No plaintext passwords stored

## Database Schema Details

### Users Table
- `id`: Auto-incrementing primary key
- `name`: User's display name
- `email`: Unique email address (used for login)
- `hashed_password`: Bcrypt-hashed password (cost 12)
- `administrator`: Boolean flag for admin privileges
- `created_at`: Timestamp of user creation

### Projects Table
- `id`: Auto-incrementing primary key
- `title`: Project name/title
- `userid`: Foreign key to users table
- `cfileids`: JSON array of associated C file IDs
- `analysis`: Binary data containing VACCS analysis results
- `created_at`: Timestamp of project creation
- `parameters`: Optional analysis parameters

### CFiles Table
- `id`: Auto-incrementing primary key
- `filename`: Name of the C source file
- `content`: Full content of the C source file
- `projectid`: Foreign key to projects table
- `created_at`: Timestamp of file upload

### Sessions Table
- `token`: Unique session token (43 characters)
- `data`: Serialized session data
- `expiry`: Session expiration timestamp with microseconds

## Testing the Installation

### 1. Verify Database Connection

```bash
mysql -u scvis -p scvis_local -e "SELECT 'Connection successful' AS status;"
```

### 2. Check Table Creation

```bash
mysql -u scvis -p scvis_local -e "SHOW TABLES;"
```

### 3. Test Application

```bash
# From project root
go run ./cmd/web
```

Visit `https://localhost:4000` to test the application.

## Security Considerations

### Database Security
- Use strong passwords for database users
- Limit database user privileges to only required databases
- Use SSL/TLS for database connections in production
- Regularly update MySQL/MariaDB

### File Permissions
- Configuration files should be readable only by application user
- Production config: `640` (root:www-data)
- Development config: `600` (user only)

### Network Security
- Restrict MySQL access to localhost only (unless clustering)
- Use firewall rules to limit database port access
- Consider using Unix sockets for local connections

## Backup and Maintenance

### Creating Backups

```bash
# Full database backup
mysqldump -u scvis -p scvis > scvis_backup_$(date +%Y%m%d).sql
mysqldump -u scvis -p scvis_local > scvis_local_backup_$(date +%Y%m%d).sql

# Structure only
mysqldump -u scvis -p --no-data scvis > scvis_structure.sql
```

### Restoring Backups

```bash
mysql -u scvis -p scvis < scvis_backup_20241201.sql
```

### Regular Maintenance

```sql
-- Optimize tables
OPTIMIZE TABLE users, projects, cfiles, sessions;

-- Check table integrity
CHECK TABLE users, projects, cfiles, sessions;

-- Clean expired sessions
DELETE FROM sessions WHERE expiry < NOW();
```

## Troubleshooting

### Common Issues

**Connection Refused:**
- Check if MySQL service is running: `sudo systemctl status mysql`
- Verify MySQL is listening: `sudo netstat -tlnp | grep 3306`

**Access Denied:**
- Verify username and password
- Check user privileges: `SHOW GRANTS FOR 'scvis'@'localhost';`

**Table Doesn't Exist:**
- Verify database selection: `USE scvis_local;`
- Check table creation: `SHOW TABLES;`

**Foreign Key Errors:**
- Ensure parent tables exist before creating child tables
- Check foreign key constraints: `SHOW CREATE TABLE projects;`

### Log Files

Check MySQL error logs:
- Ubuntu/Debian: `/var/log/mysql/error.log`
- RHEL/CentOS: `/var/log/mysqld.log`

### Performance Tuning

For production deployments, consider these MySQL settings:

```ini
[mysqld]
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
max_connections = 200
query_cache_size = 64M
```

## Support

For additional help:
- Check the main project README
- Review application logs in `/var/log/scvisual/`
- Examine MySQL error logs
- Use the user creation script in `scripts/` directory

## Integration with Build System

The database installation integrates with the existing build system:

1. **Database Setup**: Use `./install_database.sh`
2. **Build Application**: Use `./build.sh` or `./quick-build.sh`
3. **Deploy**: Use deployment scripts with systemd service

This ensures a complete development-to-production workflow.
