#!/bin/bash

# MySQL Installation and Setup Script
# This script installs MySQL server and sets up a secure root password

set -e

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

echo "================================================================="
echo "               MySQL Installation and Setup Script"
echo "================================================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    print_error "This script must be run as root (use sudo)"
    exit 1
fi

# Function to generate a secure password
generate_password() {
    openssl rand -base64 16 | tr -d "=+/" | cut -c1-16
}

# Check if MySQL is already installed
check_mysql_installed() {
    if dpkg -l | grep -q mysql-server; then
        print_warning "MySQL server is already installed"
        if systemctl is-active --quiet mysql; then
            print_info "MySQL service is running"
        else
            print_info "MySQL service is not running"
        fi
        
        # Test if root password is set
        if mysql -u root -e "SELECT 1;" &>/dev/null; then
            print_warning "MySQL root has no password set"
            return 1
        else
            print_info "MySQL root password is already configured"
            return 0
        fi
    else
        print_info "MySQL server is not installed"
        return 1
    fi
}

# Install MySQL server
install_mysql() {
    print_info "Installing MySQL server..."
    
    # Update package list
    apt-get update
    
    # Set MySQL to install without prompting for password
    export DEBIAN_FRONTEND=noninteractive
    
    # Install MySQL server
    apt-get install -y mysql-server mysql-client
    
    print_status "MySQL server installed successfully"
}

# Start and enable MySQL service
start_mysql_service() {
    print_info "Starting and enabling MySQL service..."
    
    systemctl start mysql
    systemctl enable mysql
    
    # Wait for MySQL to be ready
    sleep 3
    
    print_status "MySQL service started and enabled"
}

# Set up root password
setup_root_password() {
    print_info "Setting up MySQL root password..."
    echo ""
    
    echo "Choose password setup method:"
    echo "1) Generate a secure random password (recommended)"
    echo "2) Enter your own password"
    echo "3) Use a default password (mysql123)"
    read -p "Enter choice (1-3): " choice
    
    case $choice in
        1)
            MYSQL_ROOT_PASSWORD=$(generate_password)
            echo ""
            print_status "Generated secure password: $MYSQL_ROOT_PASSWORD"
            ;;
        2)
            echo ""
            echo -n "Enter new root password: "
            read -s MYSQL_ROOT_PASSWORD
            echo ""
            echo -n "Confirm password: "
            read -s CONFIRM_PASSWORD
            echo ""
            
            if [ "$MYSQL_ROOT_PASSWORD" != "$CONFIRM_PASSWORD" ]; then
                print_error "Passwords do not match"
                exit 1
            fi
            
            if [ -z "$MYSQL_ROOT_PASSWORD" ]; then
                print_error "Password cannot be empty"
                exit 1
            fi
            ;;
        3)
            MYSQL_ROOT_PASSWORD="mysql123"
            print_warning "Using default password: mysql123"
            print_warning "Please change this password after installation!"
            ;;
        *)
            print_error "Invalid choice"
            exit 1
            ;;
    esac
    
    echo ""
    print_info "Setting MySQL root password..."
    
    # Set the root password
    mysql -u root -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '$MYSQL_ROOT_PASSWORD';"
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "FLUSH PRIVILEGES;"
    
    print_status "MySQL root password set successfully"
    
    # Save password to secure file
    echo "$MYSQL_ROOT_PASSWORD" > /root/.mysql_root_password
    chmod 600 /root/.mysql_root_password
    print_status "Root password saved to /root/.mysql_root_password"
}

# Secure MySQL installation
secure_mysql() {
    print_info "Securing MySQL installation..."
    
    # Remove anonymous users
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "DELETE FROM mysql.user WHERE User='';"
    
    # Disable remote root login
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');"
    
    # Remove test database
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "DROP DATABASE IF EXISTS test;"
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%';"
    
    # Reload privilege tables
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "FLUSH PRIVILEGES;"
    
    print_status "MySQL installation secured"
}

# Test MySQL installation
test_mysql() {
    print_info "Testing MySQL installation..."
    
    # Test root connection
    if mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "SELECT VERSION();" &>/dev/null; then
        MYSQL_VERSION=$(mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "SELECT VERSION();" -s -N)
        print_status "MySQL connection successful - Version: $MYSQL_VERSION"
    else
        print_error "Failed to connect to MySQL"
        exit 1
    fi
    
    # Test service status
    if systemctl is-active --quiet mysql; then
        print_status "MySQL service is running"
    else
        print_error "MySQL service is not running"
        exit 1
    fi
}

# Create SCVis database user (optional)
create_scvis_user() {
    echo ""
    read -p "Create SCVis database user? [Y/n]: " CREATE_USER
    CREATE_USER=${CREATE_USER:-Y}
    
    if [[ $CREATE_USER =~ ^[Yy]$ ]]; then
        print_info "Creating SCVis database user..."
        
        # Generate password for SCVis user
        SCVIS_PASSWORD=$(generate_password)
        
        # Create databases and user
        mysql -u root -p"$MYSQL_ROOT_PASSWORD" << EOF
CREATE DATABASE IF NOT EXISTS scvis;
CREATE DATABASE IF NOT EXISTS scvis_local;
CREATE USER IF NOT EXISTS 'scvis'@'localhost' IDENTIFIED BY '$SCVIS_PASSWORD';
GRANT ALL PRIVILEGES ON scvis.* TO 'scvis'@'localhost';
GRANT ALL PRIVILEGES ON scvis_local.* TO 'scvis'@'localhost';
FLUSH PRIVILEGES;
EOF
        
        print_status "SCVis databases and user created"
        print_info "SCVis user password: $SCVIS_PASSWORD"
        
        # Save SCVis password
        echo "$SCVIS_PASSWORD" > /root/.scvis_db_password
        chmod 600 /root/.scvis_db_password
        print_status "SCVis password saved to /root/.scvis_db_password"
        
        # Create .env file
        if [ ! -f ".env" ]; then
            cat > .env << EOF
# SCVis Database Configuration
LOCAL_SERVER=true
DB_HOST=localhost
DB_USER=scvis
DB_PASSWORD=$SCVIS_PASSWORD
DB_NAME=scvis_local
PORT=8080

# Log Configuration
LOG_DIR=./log/scvisual
LOG_MAX_SIZE_MB=1
LOG_MAX_BACKUPS=5

# SSL Certificate Configuration
SSL_CERT_PATH=./tls/localhost.pem
SSL_KEY_PATH=./tls/localhost-key.pem
EOF
            chmod 600 .env
            print_status "Created .env configuration file"
        else
            print_warning ".env file already exists, not overwriting"
        fi
    fi
}

# Display summary
display_summary() {
    echo ""
    echo "================================================================="
    print_status "MySQL Installation Complete!"
    echo "================================================================="
    echo ""
    echo "Installation Summary:"
    echo "  MySQL Version: $(mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "SELECT VERSION();" -s -N)"
    echo "  Root Password: Saved to /root/.mysql_root_password"
    if [ -f "/root/.scvis_db_password" ]; then
        echo "  SCVis User: scvis"
        echo "  SCVis Password: Saved to /root/.scvis_db_password"
        echo "  Databases: scvis, scvis_local"
    fi
    echo ""
    echo "MySQL Service Commands:"
    echo "  Status: systemctl status mysql"
    echo "  Start:  systemctl start mysql"
    echo "  Stop:   systemctl stop mysql"
    echo "  Logs:   journalctl -u mysql"
    echo ""
    echo "Connect to MySQL:"
    echo "  mysql -u root -p"
    if [ -f "/root/.scvis_db_password" ]; then
        echo "  mysql -u scvis -p"
    fi
    echo ""
    echo "Next Steps:"
    echo "  1. Run: ./install_database.sh (to create SCVis tables)"
    echo "  2. Test: go run ./cmd/web"
    echo "  3. Deploy: sudo ./deploy.sh"
}

# Main installation process
main() {
    print_info "Starting MySQL installation and setup..."
    echo ""
    
    # Check if already installed
    if check_mysql_installed; then
        echo ""
        read -p "MySQL appears to be configured. Continue anyway? [y/N]: " CONTINUE
        if [[ ! $CONTINUE =~ ^[Yy]$ ]]; then
            echo "Installation cancelled."
            exit 0
        fi
    fi
    
    # Install MySQL if not present
    if ! command -v mysql &> /dev/null; then
        install_mysql
    else
        print_info "MySQL binaries already present"
    fi
    
    # Start service
    start_mysql_service
    
    # Set up root password if not already set
    if mysql -u root -e "SELECT 1;" &>/dev/null; then
        setup_root_password
        secure_mysql
    else
        print_info "Root password already configured"
        echo -n "Enter current root password to continue: "
        read -s MYSQL_ROOT_PASSWORD
        echo ""
        
        if ! mysql -u root -p"$MYSQL_ROOT_PASSWORD" -e "SELECT 1;" &>/dev/null; then
            print_error "Invalid password"
            exit 1
        fi
    fi
    
    # Test installation
    test_mysql
    
    # Optionally create SCVis user
    create_scvis_user
    
    # Display summary
    display_summary
}

# Cleanup on error
cleanup_on_error() {
    if [ $? -ne 0 ]; then
        print_error "Installation failed!"
        echo ""
        echo "To cleanup and try again:"
        echo "  sudo ./uninstall_mysql.sh"
        echo "  sudo ./install_mysql.sh"
    fi
}

# Set trap for cleanup
trap cleanup_on_error EXIT

# Run main function
main "$@"

# Remove trap on successful completion
trap - EXIT

echo ""
print_status "MySQL setup completed successfully!"
