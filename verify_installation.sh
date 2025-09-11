#!/bin/bash

# SCVis Installation Verification Script
# This script checks if all components are properly installed and configured

echo "================================================================="
echo "            SCVis Installation Verification"
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

# Check system dependencies
check_dependencies() {
    print_info "Checking system dependencies..."
    
    # Go
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | grep -o 'go[0-9.]*' | head -1)
        print_status "Go installed: $GO_VERSION"
    else
        print_error "Go is not installed"
        return 1
    fi
    
    # MySQL client
    if command -v mysql &> /dev/null; then
        print_status "MySQL client installed"
    else
        print_error "MySQL client is not installed"
        return 1
    fi
    
    # MySQL server
    if systemctl is-active --quiet mysql || systemctl is-active --quiet mysqld; then
        print_status "MySQL server is running"
    else
        print_warning "MySQL server is not running"
    fi
    
    return 0
}

# Check project structure
check_project_structure() {
    print_info "Checking project structure..."
    
    required_files=(
        "cmd/web/main.go"
        "internal/models/users.go"
        "ui/efs.go"
        "install_database.sh"
        "validate_database.sh"
        "build.sh"
        "quick-build.sh"
    )
    
    for file in "${required_files[@]}"; do
        if [ -f "$file" ]; then
            print_status "$file exists"
        else
            print_error "$file is missing"
            return 1
        fi
    done
    
    # Check if scripts are executable
    for script in install_database.sh validate_database.sh build.sh quick-build.sh; do
        if [ -x "$script" ]; then
            print_status "$script is executable"
        else
            print_warning "$script is not executable (run: chmod +x $script)"
        fi
    done
    
    return 0
}

# Check Go dependencies
check_go_dependencies() {
    print_info "Checking Go dependencies..."
    
    if [ -f "go.mod" ]; then
        print_status "go.mod exists"
        
        # Try to verify dependencies
        if go mod verify &> /dev/null; then
            print_status "Go module dependencies verified"
        else
            print_warning "Go dependencies may need updating (run: go mod tidy)"
        fi
    else
        print_error "go.mod is missing"
        return 1
    fi
    
    return 0
}

# Check database configuration
check_database() {
    print_info "Checking database configuration..."
    
    if ./validate_database.sh &> /dev/null; then
        print_status "Database is properly configured"
    else
        print_warning "Database validation failed (run: ./validate_database.sh for details)"
        return 1
    fi
    
    return 0
}

# Check build capability
check_build() {
    print_info "Testing build capability..."
    
    # Test Go build
    if go build -o /tmp/scvis-test ./cmd/web &> /dev/null; then
        print_status "Go build successful"
        rm -f /tmp/scvis-test
    else
        print_error "Go build failed"
        return 1
    fi
    
    return 0
}

# Check configuration files
check_configuration() {
    print_info "Checking configuration files..."
    
    if [ -f ".env" ]; then
        print_status "Development configuration (.env) exists"
    else
        print_warning "Development configuration (.env) missing"
    fi
    
    if [ -f "/etc/scvisual/scvisual.conf" ]; then
        print_status "Production configuration exists"
    else
        print_info "Production configuration not found (normal for development)"
    fi
    
    if [ -f "scvisual.conf.example" ]; then
        print_status "Configuration example exists"
    fi
    
    return 0
}

# Check TLS certificates
check_certificates() {
    print_info "Checking TLS certificates..."
    
    if [ -f "tls/localhost.pem" ] && [ -f "tls/localhost-key.pem" ]; then
        print_status "Development certificates exist"
    else
        print_warning "Development certificates missing (app may not start with HTTPS)"
    fi
    
    # Check for Let's Encrypt certificates (production)
    if [ -f "/etc/letsencrypt/live/*/fullchain.pem" ]; then
        print_status "Let's Encrypt certificates found"
    else
        print_info "Let's Encrypt certificates not found (normal for development)"
    fi
    
    return 0
}

# Check systemd service (if applicable)
check_service() {
    print_info "Checking systemd service..."
    
    if [ -f "scvisual.service" ]; then
        print_status "Service file exists"
        
        if [ -f "/etc/systemd/system/scvisual.service" ]; then
            print_status "Service is installed"
            
            if systemctl is-enabled scvisual &> /dev/null; then
                print_status "Service is enabled"
            else
                print_info "Service is not enabled"
            fi
            
            if systemctl is-active scvisual &> /dev/null; then
                print_status "Service is running"
            else
                print_info "Service is not running"
            fi
        else
            print_info "Service is not installed (normal for development)"
        fi
    else
        print_warning "Service file missing"
    fi
    
    return 0
}

# Generate recommendations
generate_recommendations() {
    echo ""
    echo "================================================================="
    echo "                    Recommendations"
    echo "================================================================="
    echo ""
    
    if [ $DEPS_OK -eq 0 ] && [ $STRUCTURE_OK -eq 0 ] && [ $GO_OK -eq 0 ] && [ $DB_OK -eq 0 ] && [ $BUILD_OK -eq 0 ]; then
        print_status "Your SCVis installation is ready!"
        echo ""
        echo "Quick start commands:"
        echo "  Development: ./quick-build.sh"
        echo "  Production:  ./build.sh"
        echo "  Create user: cd scripts && go run create_user.go"
        echo ""
        echo "Access the application at:"
        if [ -f ".env" ]; then
            PORT=$(grep "PORT=" .env 2>/dev/null | cut -d'=' -f2 || echo "4000")
            echo "  https://localhost:$PORT"
        else
            echo "  https://localhost:4000"
        fi
    else
        print_warning "Some issues were found. Please review the output above."
        echo ""
        echo "Common fixes:"
        echo "• Install missing dependencies"
        echo "• Run database setup: ./install_database.sh"
        echo "• Update Go dependencies: go mod tidy"
        echo "• Make scripts executable: chmod +x *.sh"
        echo ""
        echo "For detailed help, see:"
        echo "• README.md - General documentation"
        echo "• DATABASE.md - Database setup guide"
        echo "• BUILD.md - Build system documentation"
    fi
}

# Main verification process
main() {
    echo "This script will verify your SCVis installation is ready for use."
    echo ""
    
    # Initialize status variables
    DEPS_OK=1
    STRUCTURE_OK=1
    GO_OK=1
    DB_OK=1
    BUILD_OK=1
    
    # Run all checks
    if check_dependencies; then DEPS_OK=0; fi
    if check_project_structure; then STRUCTURE_OK=0; fi
    if check_go_dependencies; then GO_OK=0; fi
    if check_database; then DB_OK=0; fi
    if check_build; then BUILD_OK=0; fi
    
    # Additional checks (non-critical)
    check_configuration
    check_certificates
    check_service
    
    # Generate final recommendations
    generate_recommendations
}

# Run main function
main "$@"
