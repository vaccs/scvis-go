#!/bin/bash

# SCVis Database Installation Script
# This script creates and configures the MySQL database for the SCVis application

set -e  # Exit on any error

echo "================================================================="
echo "               SCVis Database Installation Script"
echo "================================================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Check if MySQL is installed and running
check_mysql() {
    print_info "Checking MySQL installation..."
    
    if ! command -v mysql &> /dev/null; then
        print_error "MySQL client not found. Please install MySQL."
        echo "  Ubuntu/Debian: sudo apt install mysql-client-core-8.0"
        echo "  RHEL/CentOS: sudo yum install mysql"
        exit 1
    fi
    
    if ! systemctl is-active --quiet mysql && ! systemctl is-active --quiet mysqld; then
        print_warning "MySQL service is not running. Attempting to start..."
        if systemctl start mysql 2>/dev/null || systemctl start mysqld 2>/dev/null; then
            print_status "MySQL service started successfully"
        else
            print_error "Failed to start MySQL service. Please start it manually."
            exit 1
        fi
    else
        print_status "MySQL service is running"
    fi
}

# Get MySQL root credentials
get_mysql_credentials() {
    echo ""
    print_info "Please provide MySQL root credentials to create the database:"
    
    read -p "MySQL root username [root]: " MYSQL_ROOT_USER
    MYSQL_ROOT_USER=${MYSQL_ROOT_USER:-root}
    
    echo ""
    print_info "You will be prompted for the MySQL root password during operations"
}

# Test MySQL connection
test_mysql_connection() {
    print_info "Testing MySQL connection..."
    print_info "Please enter MySQL root password when prompted:"
    
    if ! echo "SELECT 1;" | mysql -u "$MYSQL_ROOT_USER" -p; then
        print_error "Failed to connect to MySQL. Please check your credentials."
        exit 1
    fi
    
    print_status "MySQL connection successful"
}

# Get database configuration
get_database_config() {
    echo ""
    print_info "Database Configuration:"
    echo "This script will create two databases: scvis (production) and scvis_local (development)"
    echo ""
    
    read -p "Database user [scvis]: " DB_USER
    DB_USER=${DB_USER:-scvis}
    
    read -s -p "Database password for user '$DB_USER': " DB_PASSWD
    echo ""
    
    if [ -z "$DB_PASSWD" ]; then
        print_error "Database password cannot be empty"
        exit 1
    fi
    
    read -p "Production database name [scvis]: " DB_NAME
    DB_NAME=${DB_NAME:-scvis}
    
    read -p "Development database name [scvis_local]: " DB_LOCAL_NAME
    DB_LOCAL_NAME=${DB_LOCAL_NAME:-scvis_local}
    
    read -p "Database host [localhost]: " DB_HOST
    DB_HOST=${DB_HOST:-localhost}
    
    read -p "Application port [4000]: " APP_PORT
    APP_PORT=${APP_PORT:-4000}
}

# Create databases and user
create_databases() {
    print_info "Creating databases and user..."
    print_info "Please enter MySQL root password when prompted:"
    
    mysql -u "$MYSQL_ROOT_USER" -p << EOF
-- Create databases
CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
CREATE DATABASE IF NOT EXISTS \`$DB_LOCAL_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

-- Create user if it doesn't exist
CREATE USER IF NOT EXISTS '$DB_USER'@'$DB_HOST' IDENTIFIED BY '$DB_PASSWD';

-- Grant privileges
GRANT ALL PRIVILEGES ON \`$DB_NAME\`.* TO '$DB_USER'@'$DB_HOST';
GRANT ALL PRIVILEGES ON \`$DB_LOCAL_NAME\`.* TO '$DB_USER'@'$DB_HOST';

-- Flush privileges
FLUSH PRIVILEGES;
EOF
    
    print_status "Databases and user created successfully"
}

# Create table schema
create_tables() {
    local database=$1
    print_info "Creating tables in database: $database"
    print_info "Please enter password for user '$DB_USER' when prompted:"
    
    mysql -u "$DB_USER" -p "$database" << 'EOF'
-- Create users table
CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `hashed_password` char(60) NOT NULL,
  `administrator` boolean NOT NULL DEFAULT false,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_uc_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create projects table
CREATE TABLE IF NOT EXISTS `projects` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `userid` int NOT NULL,
  `cfileids` json NOT NULL,
  `analysis` blob,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `parameters` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_projects_userid` (`userid`),
  CONSTRAINT `fk_projects_userid` FOREIGN KEY (`userid`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create cfiles table
CREATE TABLE IF NOT EXISTS `cfiles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `filename` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `projectid` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_cfiles_projectid` (`projectid`),
  CONSTRAINT `fk_cfiles_projectid` FOREIGN KEY (`projectid`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create sessions table for session storage
CREATE TABLE IF NOT EXISTS `sessions` (
  `token` char(43) NOT NULL,
  `data` blob NOT NULL,
  `expiry` timestamp(6) NOT NULL,
  PRIMARY KEY (`token`),
  KEY `sessions_expiry_idx` (`expiry`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
EOF
    
    print_status "Tables created successfully in $database"
}

# Create sample admin user
create_sample_user() {
    local database=$1
    
    echo ""
    read -p "Create a sample admin user? [Y/n]: " CREATE_USER
    CREATE_USER=${CREATE_USER:-Y}
    
    if [[ $CREATE_USER =~ ^[Yy]$ ]]; then
        print_info "Creating sample admin user..."
        
        read -p "Admin username [admin]: " ADMIN_USERNAME
        ADMIN_USERNAME=${ADMIN_USERNAME:-admin}
        
        read -p "Admin email [admin@localhost]: " ADMIN_EMAIL
        ADMIN_EMAIL=${ADMIN_EMAIL:-admin@localhost}
        
        read -s -p "Admin password: " ADMIN_PASSWORD
        echo ""
        
        if [ -z "$ADMIN_PASSWORD" ]; then
            print_warning "Using default password: 'admin123'"
            ADMIN_PASSWORD="admin123"
        fi
        
        # Generate bcrypt hash (cost 12)
        HASHED_PASSWORD=$(python3 -c "import bcrypt; print(bcrypt.hashpw('$ADMIN_PASSWORD'.encode(), bcrypt.gensalt(rounds=12)).decode())" 2>/dev/null || echo "")
        
        if [ -z "$HASHED_PASSWORD" ]; then
            print_warning "Python3 with bcrypt not available. Using alternative method..."
            # Use a pre-computed hash for 'admin123' (bcrypt cost 12)
            HASHED_PASSWORD='$2b$12$G1Ct1AvS3Drg.pXzpI94FuYMTtXAZTtNQzcE2kQfb4macJwjBW3Pm'
            print_warning "Default password set to 'admin123'. Please change it after first login."
        fi
        
        print_info "Please enter password for user '$DB_USER' when prompted:"
        mysql -u "$DB_USER" -p "$database" << EOF
INSERT IGNORE INTO users (name, email, hashed_password, administrator) 
VALUES ('$ADMIN_USERNAME', '$ADMIN_EMAIL', '$HASHED_PASSWORD', true);
EOF
        
        print_status "Sample admin user created"
        echo "  Username: $ADMIN_USERNAME"
        echo "  Email: $ADMIN_EMAIL"
        if [ "$HASHED_PASSWORD" = '$2b$12$G1Ct1AvS3Drg.pXzpI94FuYMTtXAZTtNQzcE2kQfb4macJwjBW3Pm' ]; then
            echo "  Password: admin123 (please change after first login)"
        else
            echo "  Password: [as provided]"
        fi
    fi
}

# Create configuration file
create_config_file() {
    echo ""
    read -p "Create configuration file? [Y/n]: " CREATE_CONFIG
    CREATE_CONFIG=${CREATE_CONFIG:-Y}
    
    if [[ $CREATE_CONFIG =~ ^[Yy]$ ]]; then
        print_info "Creating configuration files..."
        
        # Create .env file for development
        cat > .env << EOF
DB_HOST=$DB_HOST
DB_USER=$DB_USER
DB_PASSWD=$DB_PASSWD
DB_NAME=$DB_LOCAL_NAME
PORT=$APP_PORT

# Development/Local Server Mode
LOCAL_SERVER=true

# Log Configuration
LOG_DIR=./log/scvisual
LOG_MAX_SIZE_MB=1
LOG_MAX_BACKUPS=5

# SSL Certificate Configuration
SSL_CERT_PATH=./tls/localhost.pem
SSL_KEY_PATH=./tls/localhost-key.pem
EOF
        
        # Create production .env template
        cat > production.env << EOF
# SCVis Production Configuration
# Copy this file to /opt/scvis/.env for production deployment

DB_HOST=$DB_HOST
DB_USER=$DB_USER
DB_PASSWD=$DB_PASSWD
DB_NAME=$DB_NAME
PORT=$APP_PORT

# Production Server Mode
LOCAL_SERVER=false

# Log Configuration
LOG_DIR=/var/log/scvisual
LOG_MAX_SIZE_MB=10
LOG_MAX_BACKUPS=5

# SSL Certificate Configuration
SSL_CERT_PATH=/etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem
SSL_KEY_PATH=/etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem
EOF
        
        chmod 600 .env production.env
        
        print_status "Configuration files created:"
        echo "  .env (development) - using database: $DB_LOCAL_NAME"
        echo "  production.env (production template) - using database: $DB_NAME"
        echo ""
        print_warning "For production deployment:"
        echo "  sudo mkdir -p /opt/scvis"
        echo "  sudo cp production.env /opt/scvis/.env"
        echo "  sudo chown root:www-data /opt/scvis/.env"
        echo "  sudo chmod 640 /opt/scvis/.env"
        echo ""
        print_info "Legacy /etc/scvisual/scvisual.conf is no longer used"
    fi
}

# Verify installation
verify_installation() {
    print_info "Verifying installation..."
    print_info "Please enter password for user '$DB_USER' when prompted for verification:"
    
    # Test database connections
    for db in "$DB_NAME" "$DB_LOCAL_NAME"; do
        if echo "SELECT COUNT(*) FROM users;" | mysql -u "$DB_USER" -p "$db" &> /dev/null; then
            print_status "Database $db is accessible and functional"
        else
            print_error "Failed to verify database $db"
        fi
    done
    
    # Test table creation
    table_count=$(echo "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='$DB_LOCAL_NAME';" | mysql -u "$DB_USER" -p -s 2>/dev/null || echo "0")
    
    if [ "$table_count" -ge "4" ]; then
        print_status "All required tables created successfully"
    else
        print_warning "Some tables may be missing (found $table_count tables, expected 4)"
    fi
}

# Cleanup on error
cleanup_on_error() {
    if [ $? -ne 0 ]; then
        print_error "Installation failed. You may need to manually clean up:"
        echo "  DROP DATABASE IF EXISTS \`$DB_NAME\`;"
        echo "  DROP DATABASE IF EXISTS \`$DB_LOCAL_NAME\`;"
        echo "  DROP USER IF EXISTS '$DB_USER'@'$DB_HOST';"
    fi
}

# Set trap for cleanup
trap cleanup_on_error EXIT

# Main installation process
main() {
    echo "This script will:"
    echo "• Check MySQL installation and service status"
    echo "• Create production and development databases"
    echo "• Create database user with appropriate permissions"
    echo "• Create all required tables with proper relationships"
    echo "• Optionally create a sample admin user"
    echo "• Generate configuration files"
    echo ""
    print_warning "Password Security Note:"
    echo "This script will prompt for passwords interactively (not on command line)"
    echo "You will be asked for:"
    echo "  1. MySQL root password (multiple times for different operations)"
    echo "  2. New database user password (for verification steps)"
    echo ""
    
    read -p "Continue with installation? [Y/n]: " CONTINUE
    CONTINUE=${CONTINUE:-Y}
    
    if [[ ! $CONTINUE =~ ^[Yy]$ ]]; then
        echo "Installation cancelled."
        exit 0
    fi
    
    echo ""
    echo "Starting SCVis database installation..."
    echo "======================================="
    
    # Run installation steps
    check_mysql
    get_mysql_credentials
    test_mysql_connection
    get_database_config
    create_databases
    create_tables "$DB_NAME"
    create_tables "$DB_LOCAL_NAME"
    create_sample_user "$DB_LOCAL_NAME"
    create_config_file
    verify_installation
    
    echo ""
    echo "================================================================="
    print_status "SCVis Database Installation Complete!"
    echo "================================================================="
    echo ""
    echo "Installation Summary:"
    echo "  Production Database: $DB_NAME"
    echo "  Development Database: $DB_LOCAL_NAME"
    echo "  Database User: $DB_USER"
    echo "  Database Host: $DB_HOST"
    echo "  Application Port: $APP_PORT"
    echo ""
    echo "Next Steps:"
    echo "1. Test the application with: go run ./cmd/web"
    echo "2. For production, copy production.env to /opt/scvis/.env"
    echo "3. Build and deploy the application using the build scripts"
    echo ""
    echo "For additional users, run:"
    echo "  cd scripts && go run create_user.go"
}

# Remove trap on successful completion
trap - EXIT

# Run main function
main "$@"
