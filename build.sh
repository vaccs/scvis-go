#!/bin/bash

# Comprehensive Build Script for SCVis Go Application
# This script builds the Flutter web frontend and the Go backend

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "\n${BLUE}===== $1 =====${NC}"
}

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$SCRIPT_DIR"
FLUTTER_DIR="$SCRIPT_DIR/../scvis-flutter"
FLUTTER_PROJECT_DIR="$FLUTTER_DIR/scvis_flutter"
FLUTTER_BUILD_DIR="$FLUTTER_PROJECT_DIR/build/web"
GO_VISUAL_DIR="$PROJECT_DIR/visual/scvis"

print_header "SCVis Full Build Process"
print_status "Project directory: $PROJECT_DIR"
print_status "Flutter directory: $FLUTTER_DIR"
print_status "Target visual directory: $GO_VISUAL_DIR"

# Step 1: Build Flutter Web Application
print_header "Step 1: Building Flutter Web Application"

# Check if Flutter project exists
if [ ! -d "$FLUTTER_DIR" ]; then
    print_error "Flutter project directory not found: $FLUTTER_DIR"
    exit 1
fi

if [ ! -f "$FLUTTER_DIR/build_web.sh" ]; then
    print_error "build_web.sh not found in: $FLUTTER_DIR"
    exit 1
fi

print_status "Changing to Flutter directory: $FLUTTER_DIR"
cd "$FLUTTER_DIR"

print_status "Running Flutter build script..."
if ./build_web.sh; then
    print_success "Flutter web build completed successfully"
else
    print_error "Flutter web build failed"
    exit 1
fi

# Step 2: Copy Flutter build output to Go project
print_header "Step 2: Copying Flutter Build to Go Project"

# Check if Flutter build output exists
if [ ! -d "$FLUTTER_BUILD_DIR" ]; then
    print_error "Flutter build output not found: $FLUTTER_BUILD_DIR"
    exit 1
fi

print_status "Source: $FLUTTER_BUILD_DIR"
print_status "Destination: $GO_VISUAL_DIR"

# Create visual directory if it doesn't exist
mkdir -p "$GO_VISUAL_DIR"

# Remove old files but preserve any custom files
print_status "Cleaning destination directory..."
if [ -d "$GO_VISUAL_DIR" ]; then
    # Create a backup of any custom files that might exist
    if [ -f "$GO_VISUAL_DIR/.gitkeep" ]; then
        cp "$GO_VISUAL_DIR/.gitkeep" "/tmp/.gitkeep.backup" 2>/dev/null || true
    fi
    
    # Remove old build files
    rm -rf "$GO_VISUAL_DIR"/*
    
    # Restore custom files
    if [ -f "/tmp/.gitkeep.backup" ]; then
        cp "/tmp/.gitkeep.backup" "$GO_VISUAL_DIR/.gitkeep"
        rm "/tmp/.gitkeep.backup"
    fi
fi

# Copy new Flutter build
print_status "Copying Flutter build files..."
cp -r "$FLUTTER_BUILD_DIR"/* "$GO_VISUAL_DIR/"

# Verify copy was successful
if [ -f "$GO_VISUAL_DIR/index.html" ]; then
    print_success "Flutter files copied successfully"
    
    # Show some stats
    file_count=$(find "$GO_VISUAL_DIR" -type f | wc -l)
    total_size=$(du -sh "$GO_VISUAL_DIR" | cut -f1)
    print_status "Copied $file_count files, total size: $total_size"
else
    print_error "Failed to copy Flutter files - index.html not found"
    exit 1
fi

# Step 3: Build Go Application
print_header "Step 3: Building Go Application"

# Change back to Go project directory
cd "$PROJECT_DIR"

print_status "Current directory: $(pwd)"

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found. Are you in the correct Go project directory?"
    exit 1
fi

# Clean previous Go build
print_status "Cleaning previous Go builds..."
rm -f scvisual

# Build Go application
print_status "Building Go application..."
if go build -o scvisual ./cmd/web; then
    print_success "Go application built successfully"
    
    # Show binary info
    if [ -f "scvisual" ]; then
        binary_size=$(du -sh scvisual | cut -f1)
        print_status "Binary size: $binary_size"
        print_status "Binary location: $PROJECT_DIR/scvisual"
    fi
else
    print_error "Go build failed"
    exit 1
fi

# Final verification
print_header "Build Verification"

# Check Go binary
if [ -x "scvisual" ]; then
    print_success "✓ Go binary is executable"
else
    print_warning "⚠ Go binary is not executable or not found"
fi

# Check Flutter files
if [ -f "$GO_VISUAL_DIR/index.html" ] && [ -f "$GO_VISUAL_DIR/main.dart.js" ]; then
    print_success "✓ Flutter web files are present"
else
    print_warning "⚠ Flutter web files may be missing"
fi

# Build summary
print_header "Build Summary"
echo -e "${GREEN}🎉 Build completed successfully!${NC}"
echo ""
echo "📋 What was built:"
echo "   1. ✅ Flutter web application (from $FLUTTER_DIR)"
echo "   2. ✅ Web files copied to visual/scvis/"
echo "   3. ✅ Go application binary: scvisual"
echo ""
echo "🚀 Next steps:"
echo "   • Test locally: ./scvisual"
echo "   • Deploy: sudo ./deploy.sh"
echo "   • Check status: sudo systemctl status scvisual"
echo ""
echo "📁 Key locations:"
echo "   • Go binary: $PROJECT_DIR/scvisual"
echo "   • Web files: $GO_VISUAL_DIR/"
echo "   • Flutter source: $FLUTTER_DIR"
echo ""

# Optional: Show recent build info
if [ -f "$GO_VISUAL_DIR/.last_build_id" ]; then
    build_id=$(cat "$GO_VISUAL_DIR/.last_build_id" 2>/dev/null || echo "unknown")
    echo "📝 Build ID: $build_id"
fi

print_success "Full build process completed in $(pwd)"
