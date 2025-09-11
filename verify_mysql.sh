#!/bin/bash

# Quick MySQL status check script
# Use this to verify your MySQL installation

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

echo "MySQL Installation Verification"
echo "==============================="
echo ""

# Check if MySQL is installed
if command -v mysql &> /dev/null; then
    print_status "MySQL client is installed"
else
    print_error "MySQL client is not installed"
    exit 1
fi

# Check if MySQL server is installed
if dpkg -l | grep -q mysql-server; then
    print_status "MySQL server is installed"
else
    print_error "MySQL server is not installed"
    exit 1
fi

# Check if MySQL service is running
if systemctl is-active --quiet mysql; then
    print_status "MySQL service is running"
else
    print_error "MySQL service is not running"
    echo "  Try: sudo systemctl start mysql"
    exit 1
fi

# Check if root password file exists
if [ -f "/root/.mysql_root_password" ]; then
    print_status "Root password file exists"
    ROOT_PASSWORD=$(sudo cat /root/.mysql_root_password)
else
    print_warning "Root password file not found"
    echo -n "Enter MySQL root password: "
    read -s ROOT_PASSWORD
    echo ""
fi

# Test root connection
if mysql -u root -p"$ROOT_PASSWORD" -e "SELECT 1;" &>/dev/null; then
    print_status "Root connection successful"
    VERSION=$(mysql -u root -p"$ROOT_PASSWORD" -e "SELECT VERSION();" -s -N)
    echo "  Version: $VERSION"
else
    print_error "Failed to connect as root"
    exit 1
fi

# Check for SCVis databases
print_info "Checking for SCVis databases..."
DATABASES=$(mysql -u root -p"$ROOT_PASSWORD" -e "SHOW DATABASES LIKE 'scvis%';" -s -N)
if echo "$DATABASES" | grep -q "scvis"; then
    print_status "SCVis databases found:"
    echo "$DATABASES" | sed 's/^/  - /'
else
    print_warning "No SCVis databases found"
fi

# Check for SCVis user
if mysql -u root -p"$ROOT_PASSWORD" -e "SELECT User FROM mysql.user WHERE User='scvis';" -s -N | grep -q "scvis"; then
    print_status "SCVis user exists"
    
    # Test SCVis user connection
    if [ -f "/root/.scvis_db_password" ]; then
        SCVIS_PASSWORD=$(sudo cat /root/.scvis_db_password)
        if mysql -u scvis -p"$SCVIS_PASSWORD" -e "SELECT 1;" &>/dev/null; then
            print_status "SCVis user connection successful"
        else
            print_error "SCVis user connection failed"
        fi
    else
        print_warning "SCVis password file not found"
    fi
else
    print_warning "SCVis user not found"
fi

# Check configuration files
echo ""
print_info "Checking configuration files..."
if [ -f ".env" ]; then
    print_status ".env file exists"
    if grep -q "DB_PASSWORD=" .env; then
        print_status "Database password configured in .env"
    else
        print_warning "No database password in .env"
    fi
else
    print_warning ".env file not found"
fi

echo ""
print_status "MySQL verification complete!"
echo ""
echo "Connection Commands:"
echo "  Root: mysql -u root -p"
if [ -f "/root/.scvis_db_password" ]; then
    echo "  SCVis: mysql -u scvis -p"
fi
echo ""
echo "Service Commands:"
echo "  Status: sudo systemctl status mysql"
echo "  Logs: sudo journalctl -u mysql -f"
