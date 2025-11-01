#!/bin/bash

# MySQL Complete Uninstall Script
# This script completely removes MySQL server and all associated data
# WARNING: This will permanently delete all MySQL databases and data!

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
echo "               MySQL Complete Uninstall Script"
echo "================================================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    print_error "This script must be run as root (use sudo)"
    exit 1
fi

# Display warning
display_warning() {
    echo "⚠️  WARNING: This will completely remove MySQL and ALL databases!"
    echo "⚠️  All data will be permanently lost and cannot be recovered!"
    echo ""
    echo "This includes:"
    echo "  • All MySQL packages (server, client, common)"
    echo "  • All databases (including scvis and scvis_local)"
    echo "  • All MySQL configuration files"
    echo "  • All MySQL data directories"
    echo "  • MySQL user and group accounts"
    echo "  • All MySQL log files"
    echo "  • SCVis database passwords and configuration"
    echo ""
}

# Check what MySQL components are installed
check_mysql_components() {
    print_info "Checking MySQL installation..."
    
    MYSQL_INSTALLED=false
    MYSQL_RUNNING=false
    MYSQL_ENABLED=false
    SCVIS_CONFIGURED=false
    
    # Check if MySQL packages are installed
    if dpkg -l | grep -q mysql-server; then
        MYSQL_INSTALLED=true
        print_warning "MySQL server is installed"
    else
        print_info "MySQL server is not installed"
    fi
    
    # Check if MySQL service is running
    if systemctl is-active --quiet mysql 2>/dev/null; then
        MYSQL_RUNNING=true
        print_warning "MySQL service is running"
    else
        print_info "MySQL service is not running"
    fi
    
    # Check if MySQL service is enabled
    if systemctl is-enabled --quiet mysql 2>/dev/null; then
        MYSQL_ENABLED=true
        print_warning "MySQL service is enabled"
    else
        print_info "MySQL service is not enabled"
    fi
    
    # Check for SCVis configuration
    if [ -f "/root/.scvis_db_password" ] || [ -f ".env" ]; then
        SCVIS_CONFIGURED=true
        print_warning "SCVis database configuration found"
    else
        print_info "No SCVis configuration found"
    fi
    
    echo ""
    
    # If nothing is installed, exit
    if [ "$MYSQL_INSTALLED" = false ] && [ "$SCVIS_CONFIGURED" = false ]; then
        print_status "No MySQL installation found. Nothing to uninstall."
        exit 0
    fi
}

# Get confirmation
get_confirmation() {
    display_warning
    
    echo "Type 'YES I UNDERSTAND' to proceed with complete removal:"
    read -p "> " confirmation
    
    if [ "$confirmation" != "YES I UNDERSTAND" ]; then
        echo ""
        print_info "Uninstall cancelled by user."
        exit 0
    fi
    
    echo ""
    print_info "Proceeding with MySQL uninstall..."
}

# Backup databases (optional)
offer_backup() {
    if [ "$MYSQL_INSTALLED" = true ] && [ "$MYSQL_RUNNING" = true ]; then
        echo ""
        read -p "Create a backup of databases before removal? [y/N]: " backup_choice
        
        if [[ $backup_choice =~ ^[Yy]$ ]]; then
            print_info "Creating database backup..."
            
            BACKUP_DIR="mysql_backup_$(date +%Y%m%d_%H%M%S)"
            mkdir -p "$BACKUP_DIR"
            
            # Try to get root password
            ROOT_PASSWORD=""
            if [ -f "/root/.mysql_root_password" ]; then
                ROOT_PASSWORD=$(cat /root/.mysql_root_password)
                print_info "Using saved root password"
            else
                echo -n "Enter MySQL root password (or press Enter to skip backup): "
                read -s ROOT_PASSWORD
                echo ""
            fi
            
            if [ -n "$ROOT_PASSWORD" ]; then
                # Backup all databases
                if mysql -u root -p"$ROOT_PASSWORD" -e "SHOW DATABASES;" &>/dev/null; then
                    mysqldump -u root -p"$ROOT_PASSWORD" --all-databases > "$BACKUP_DIR/all_databases.sql"
                    
                    # Individual backups for SCVis databases
                    if mysql -u root -p"$ROOT_PASSWORD" -e "USE scvis;" &>/dev/null; then
                        mysqldump -u root -p"$ROOT_PASSWORD" scvis > "$BACKUP_DIR/scvis.sql"
                    fi
                    
                    if mysql -u root -p"$ROOT_PASSWORD" -e "USE scvis_local;" &>/dev/null; then
                        mysqldump -u root -p"$ROOT_PASSWORD" scvis_local > "$BACKUP_DIR/scvis_local.sql"
                    fi
                    
                    # Copy configuration files
                    cp /root/.mysql_root_password "$BACKUP_DIR/" 2>/dev/null || true
                    cp /root/.scvis_db_password "$BACKUP_DIR/" 2>/dev/null || true
                    cp .env "$BACKUP_DIR/" 2>/dev/null || true
                    
                    print_status "Backup created in: $BACKUP_DIR"
                else
                    print_error "Failed to connect to MySQL for backup"
                fi
            else
                print_warning "Skipping backup due to missing password"
            fi
        fi
    fi
}

# Stop MySQL services
stop_mysql_services() {
    if [ "$MYSQL_RUNNING" = true ]; then
        print_info "Stopping MySQL services..."
        
        systemctl stop mysql 2>/dev/null || true
        systemctl stop mysqld 2>/dev/null || true
        
        # Wait for services to stop
        sleep 3
        
        print_status "MySQL services stopped"
    fi
}

# Disable MySQL services
disable_mysql_services() {
    if [ "$MYSQL_ENABLED" = true ]; then
        print_info "Disabling MySQL services..."
        
        systemctl disable mysql 2>/dev/null || true
        systemctl disable mysqld 2>/dev/null || true
        
        print_status "MySQL services disabled"
    fi
}

# Remove MySQL packages
remove_mysql_packages() {
    if [ "$MYSQL_INSTALLED" = true ]; then
        print_info "Removing MySQL packages..."
        
        # Remove all MySQL packages
        apt-get remove --purge -y \
            mysql-server \
            mysql-client \
            mysql-common \
            mysql-server-core-* \
            mysql-client-core-* \
            mysql-server-8.0 \
            mysql-client-8.0 \
            libmysqlclient* 2>/dev/null || true
        
        # Clean up package dependencies
        apt-get autoremove -y
        apt-get autoclean
        
        print_status "MySQL packages removed"
    fi
}

# Remove MySQL data directories
remove_mysql_data() {
    print_info "Removing MySQL data directories..."
    
    # Remove main data directories
    rm -rf /var/lib/mysql
    rm -rf /var/lib/mysql-files
    rm -rf /var/lib/mysql-keyring
    rm -rf /var/log/mysql
    rm -rf /etc/mysql
    rm -rf /run/mysqld
    rm -rf /var/run/mysqld
    
    print_status "MySQL data directories removed"
}

# Remove MySQL user and group
remove_mysql_user() {
    print_info "Removing MySQL user and group..."
    
    # Remove MySQL user
    if id "mysql" &>/dev/null; then
        userdel mysql 2>/dev/null || true
        print_status "MySQL user removed"
    else
        print_info "MySQL user does not exist"
    fi
    
    # Remove MySQL group
    if getent group mysql &>/dev/null; then
        groupdel mysql 2>/dev/null || true
        print_status "MySQL group removed"
    else
        print_info "MySQL group does not exist"
    fi
}

# Remove configuration and password files
remove_config_files() {
    print_info "Removing configuration and password files..."
    
    # Remove MySQL configuration files
    rm -rf /etc/mysql*
    rm -rf /etc/apparmor.d/usr.sbin.mysqld
    
    # Remove password files
    rm -f /root/.mysql_root_password
    rm -f /root/.scvis_db_password
    rm -f /root/.mysql_history
    
    # Remove package configuration files
    rm -rf /var/lib/dpkg/info/mysql*
    
    print_status "Configuration files removed"
}

# Clean up SCVis configuration
cleanup_scvis_config() {
    if [ "$SCVIS_CONFIGURED" = true ]; then
        print_info "Cleaning up SCVis configuration..."
        
        # Backup and clean .env file
        if [ -f ".env" ]; then
            cp .env .env.backup.$(date +%Y%m%d_%H%M%S)
            
            # Comment out database-related lines
            sed -i 's/^DB_/#DB_/' .env
            sed -i 's/^LOCAL_SERVER/#LOCAL_SERVER/' .env
            
            print_status "SCVis .env file updated (database settings commented out)"
        fi
        
        # Clean production configuration if it exists
        if [ -f "/opt/scvisual/.env" ]; then
            cp /opt/scvisual/.env /opt/scvisual/.env.backup.$(date +%Y%m%d_%H%M%S)
            sed -i 's/^DB_/#DB_/' /opt/scvisual/.env
            print_status "Production .env file updated"
        fi
    fi
}

# Clean package system
clean_package_system() {
    print_info "Cleaning package system..."
    
    # Update package lists
    apt-get update
    
    # Fix any broken packages
    dpkg --configure -a
    
    # Clean package cache
    apt-get clean
    
    print_status "Package system cleaned"
}

# Verify removal
verify_removal() {
    print_info "Verifying MySQL removal..."
    
    ISSUES_FOUND=false
    
    # Check for remaining packages
    if dpkg -l | grep -q mysql; then
        print_warning "Some MySQL packages may still be installed:"
        dpkg -l | grep mysql | awk '{print "  " $2}'
        ISSUES_FOUND=true
    else
        print_status "All MySQL packages removed"
    fi
    
    # Check for remaining processes
    if pgrep -f mysql &>/dev/null; then
        print_warning "MySQL processes still running:"
        pgrep -f mysql | xargs ps -p
        ISSUES_FOUND=true
    else
        print_status "No MySQL processes running"
    fi
    
    # Check for remaining directories
    REMAINING_DIRS=""
    for dir in /var/lib/mysql /var/log/mysql /etc/mysql /run/mysqld; do
        if [ -d "$dir" ]; then
            REMAINING_DIRS="$REMAINING_DIRS $dir"
        fi
    done
    
    if [ -n "$REMAINING_DIRS" ]; then
        print_warning "Some directories still exist:$REMAINING_DIRS"
        ISSUES_FOUND=true
    else
        print_status "All MySQL directories removed"
    fi
    
    # Check for MySQL user
    if id "mysql" &>/dev/null; then
        print_warning "MySQL user still exists"
        ISSUES_FOUND=true
    else
        print_status "MySQL user removed"
    fi
    
    if [ "$ISSUES_FOUND" = false ]; then
        print_status "MySQL completely removed"
    else
        print_warning "Some components may require manual cleanup"
    fi
}

# Display summary
display_summary() {
    echo ""
    echo "================================================================="
    print_status "MySQL Uninstall Complete!"
    echo "================================================================="
    echo ""
    echo "Removed Components:"
    echo "  ✓ MySQL server and client packages"
    echo "  ✓ All databases and data"
    echo "  ✓ Configuration files"
    echo "  ✓ MySQL user and group"
    echo "  ✓ Log files and data directories"
    echo "  ✓ SCVis database configuration"
    echo ""
    
    if [ -d "mysql_backup_"* ] 2>/dev/null; then
        echo "Backup Location:"
        ls -d mysql_backup_* 2>/dev/null | head -1 | sed 's/^/  /'
        echo ""
    fi
    
    echo "To reinstall MySQL:"
    echo "  sudo ./install_mysql.sh"
    echo ""
    echo "To restore from backup (if created):"
    echo "  sudo ./install_mysql.sh"
    echo "  mysql -u root -p < backup_directory/all_databases.sql"
    echo ""
    echo "SCVis Application:"
    echo "  ✓ Configuration files backed up"
    echo "  ✓ Database settings commented out in .env"
    echo "  ! Application will not start without database"
}

# Main uninstall process
main() {
    print_info "Starting MySQL uninstall process..."
    echo ""
    
    # Check what's installed
    check_mysql_components
    
    # Get user confirmation
    get_confirmation
    
    # Offer backup
    offer_backup
    
    # Stop services
    stop_mysql_services
    disable_mysql_services
    
    # Remove MySQL
    remove_mysql_packages
    remove_mysql_data
    remove_mysql_user
    remove_config_files
    
    # Clean up SCVis
    cleanup_scvis_config
    
    # Clean package system
    clean_package_system
    
    # Verify removal
    verify_removal
    
    # Display summary
    display_summary
}

# Error handling
cleanup_on_error() {
    if [ $? -ne 0 ]; then
        print_error "Uninstall process encountered an error!"
        echo ""
        echo "You may need to manually clean up remaining components."
        echo "Check the output above for specific issues."
    fi
}

# Set trap for cleanup
trap cleanup_on_error EXIT

# Run main function
main "$@"

# Remove trap on successful completion
trap - EXIT

echo ""
print_status "MySQL uninstall completed successfully!"
