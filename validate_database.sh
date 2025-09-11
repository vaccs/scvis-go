#!/bin/bash

# SCVis Database Validation Script
# This script validates that the SCVis database is properly configured

set -e

echo "================================================================="
echo "               SCVis Database Validation Script"
echo "================================================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Load configuration
load_config() {
    if [ -f ".env" ]; then
        source .env
        print_status "Loaded development configuration from .env"
    elif [ -f "/etc/scvisual/scvisual.conf" ]; then
        source /etc/scvisual/scvisual.conf
        print_status "Loaded production configuration from /etc/scvisual/scvisual.conf"
    else
        print_error "No configuration file found (.env or /etc/scvisual/scvisual.conf)"
        exit 1
    fi
    
    # Set defaults
    DB_HOST=${DB_HOST:-localhost}
    
    # Validate required variables
    if [ -z "$DB_USER" ] || [ -z "$DB_PASSWD" ] || [ -z "$DB_NAME" ]; then
        print_error "Missing required database configuration (DB_USER, DB_PASSWD, DB_NAME)"
        exit 1
    fi
}

# Test MySQL connection
test_connection() {
    print_info "Testing database connection..."
    
    if mysql -u "$DB_USER" -p"$DB_PASSWD" -h "$DB_HOST" "$DB_NAME" -e "SELECT 1;" &> /dev/null; then
        print_status "Database connection successful"
    else
        print_error "Failed to connect to database $DB_NAME"
        return 1
    fi
}

# Check database exists
check_database() {
    print_info "Checking database existence..."
    
    if mysql -u "$DB_USER" -p"$DB_PASSWD" -h "$DB_HOST" -e "USE $DB_NAME;" &> /dev/null; then
        print_status "Database $DB_NAME exists and is accessible"
    else
        print_error "Database $DB_NAME does not exist or is not accessible"
        return 1
    fi
}

# Check required tables
check_tables() {
    print_info "Checking required tables..."
    
    required_tables=("users" "projects" "cfiles" "sessions")
    
    for table in "${required_tables[@]}"; do
        if mysql -u "$DB_USER" -p"$DB_PASSWD" -h "$DB_HOST" "$DB_NAME" -e "DESCRIBE $table;" &> /dev/null; then
            print_status "Table $table exists"
        else
            print_error "Table $table does not exist"
            return 1
        fi
    done
}

# Check table structure
check_table_structure() {
    print_info "Validating table structures..."
    
    # Check users table structure
    users_columns=$(mysql -u "$DB_USER" -p"$DB_PASSWD" -h "$DB_HOST" "$DB_NAME" -e "SHOW COLUMNS FROM users;" 2>/dev/null | wc -l)
    if [ "$users_columns" -ge "5" ]; then
        print_status "Users table structure is correct"
    else
        print_warning "Users table may have incorrect structure"
    fi
    
    # Check foreign key constraints
    if mysql -u "$DB_USER" -p"$DB_PASSWD" -h "$DB_HOST" "$DB_NAME" -e "SELECT * FROM information_schema.key_column_usage WHERE table_name='projects' AND constraint_name='fk_projects_userid';" 2>/dev/null | grep -q "fk_projects_userid"; then
        print_status "Foreign key constraints exist"
    else
        print_warning "Foreign key constraints may be missing"
    fi
}

# Check for test data
check_test_data() {
    print_info "Checking for test data..."
    
    user_count=$(mysql -u "$DB_USER" -p"$DB_PASSWD" -h "$DB_HOST" "$DB_NAME" -e "SELECT COUNT(*) FROM users;" -s 2>/dev/null || echo "0")
    
    if [ "$user_count" -gt "0" ]; then
        print_status "Found $user_count user(s) in database"
        
        # Check for admin user
        admin_count=$(mysql -u "$DB_USER" -p"$DB_PASSWD" -h "$DB_HOST" "$DB_NAME" -e "SELECT COUNT(*) FROM users WHERE administrator=1;" -s 2>/dev/null || echo "0")
        if [ "$admin_count" -gt "0" ]; then
            print_status "Found $admin_count administrator(s)"
        else
            print_warning "No administrator users found"
        fi
    else
        print_warning "No users found in database"
    fi
}

# Test application connectivity
test_app_connectivity() {
    print_info "Testing application integration..."
    
    # Try to run a quick database query using Go
    if command -v go &> /dev/null; then
        cd "$(dirname "$0")"
        
        cat > /tmp/scvis_db_test.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    err := godotenv.Load(".env")
    if err != nil {
        err = godotenv.Load("/etc/scvisual/scvisual.conf")
        if err != nil {
            fmt.Println("Failed to load configuration")
            os.Exit(1)
        }
    }

    dbUser := os.Getenv("DB_USER")
    dbPasswd := os.Getenv("DB_PASSWD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    
    if dbHost == "" {
        dbHost = "localhost"
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPasswd, dbHost, dbName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        fmt.Printf("Failed to connect: %v\n", err)
        os.Exit(1)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        fmt.Printf("Failed to ping: %v\n", err)
        os.Exit(1)
    }

    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    if err != nil {
        fmt.Printf("Failed to query: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("SUCCESS: Found %d users\n", count)
}
EOF
        
        if go run /tmp/scvis_db_test.go 2>/dev/null; then
            print_status "Go application can connect to database"
        else
            print_warning "Go application connection test failed (may be due to missing dependencies)"
        fi
        
        rm -f /tmp/scvis_db_test.go
    else
        print_warning "Go not available, skipping application connectivity test"
    fi
}

# Generate summary report
generate_report() {
    echo ""
    echo "================================================================="
    echo "                    Validation Summary"
    echo "================================================================="
    echo ""
    echo "Database Configuration:"
    echo "  Host: $DB_HOST"
    echo "  User: $DB_USER"
    echo "  Database: $DB_NAME"
    echo "  Port: ${PORT:-4000}"
    echo ""
    
    # Overall status
    if $CONNECTION_OK && $DATABASE_OK && $TABLES_OK; then
        print_status "Database is properly configured and ready for use"
        echo ""
        echo "Next steps:"
        echo "• Start the application: go run ./cmd/web"
        echo "• Access the web interface: https://localhost:${PORT:-4000}"
        echo "• Create additional users: cd scripts && go run create_user.go"
    else
        print_error "Database configuration has issues that need to be resolved"
        echo ""
        echo "Troubleshooting:"
        echo "• Run the installation script: ./install_database.sh"
        echo "• Check MySQL service: sudo systemctl status mysql"
        echo "• Verify credentials in configuration file"
    fi
}

# Main validation process
main() {
    echo "This script will validate your SCVis database configuration."
    echo ""
    
    # Initialize status variables
    CONNECTION_OK=false
    DATABASE_OK=false
    TABLES_OK=false
    
    # Load configuration
    load_config
    
    # Run validation tests
    if test_connection; then
        CONNECTION_OK=true
        
        if check_database; then
            DATABASE_OK=true
            
            if check_tables; then
                TABLES_OK=true
                check_table_structure
                check_test_data
                test_app_connectivity
            fi
        fi
    fi
    
    # Generate final report
    generate_report
}

# Run main function
main "$@"
