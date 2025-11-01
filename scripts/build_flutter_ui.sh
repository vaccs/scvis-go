#!/bin/bash

# Flutter Web UI Build Script for SCVis
# This script builds the Flutter web application and copies it to the visual/scvis directory

set -e

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

echo "========================================"
echo "    SCVis Flutter Web UI Builder"
echo "========================================"
echo ""

# Configuration
FLUTTER_PROJECT_DIR="../scvis-flutter"
BUILD_SCRIPT="build_web.sh"
SOURCE_DIR="${FLUTTER_PROJECT_DIR}/scvis_flutter/build/web"
TARGET_DIR="./visual/scvis"
CURRENT_DIR=$(pwd)

print_info "Current directory: $CURRENT_DIR"
print_info "Flutter project directory: $FLUTTER_PROJECT_DIR"
print_info "Target directory: $TARGET_DIR"

# Check if Flutter project directory exists
if [ ! -d "$FLUTTER_PROJECT_DIR" ]; then
    print_error "Flutter project directory not found: $FLUTTER_PROJECT_DIR"
    echo "Please ensure the scvis-flutter project is located at: $FLUTTER_PROJECT_DIR"
    exit 1
fi

print_status "Flutter project directory found"

# Check if build script exists
if [ ! -f "$FLUTTER_PROJECT_DIR/$BUILD_SCRIPT" ]; then
    print_error "Build script not found: $FLUTTER_PROJECT_DIR/$BUILD_SCRIPT"
    echo "Please ensure the build_web.sh script exists in the Flutter project"
    exit 1
fi

print_status "Build script found"

# Make sure the build script is executable
if [ ! -x "$FLUTTER_PROJECT_DIR/$BUILD_SCRIPT" ]; then
    print_info "Making build script executable..."
    chmod +x "$FLUTTER_PROJECT_DIR/$BUILD_SCRIPT"
fi

# Change to Flutter project directory
print_info "Changing to Flutter project directory..."
cd "$FLUTTER_PROJECT_DIR"

# Run the build script
print_info "Building Flutter web application..."
echo "Running: ./$BUILD_SCRIPT"
if ./$BUILD_SCRIPT; then
    print_status "Flutter web build completed successfully"
else
    print_error "Flutter web build failed"
    cd "$CURRENT_DIR"
    exit 1
fi

# Check if build output exists
if [ ! -d "$SOURCE_DIR" ]; then
    print_error "Build output directory not found: $SOURCE_DIR"
    print_error "The build may have failed or the output directory is different"
    cd "$CURRENT_DIR"
    exit 1
fi

print_status "Build output directory found"

# Return to original directory
cd "$CURRENT_DIR"

# Create target directory if it doesn't exist
if [ ! -d "$TARGET_DIR" ]; then
    print_info "Creating target directory: $TARGET_DIR"
    mkdir -p "$TARGET_DIR"
fi

# Backup existing files if they exist
if [ "$(ls -A $TARGET_DIR 2>/dev/null)" ]; then
    BACKUP_DIR="${TARGET_DIR}_backup_$(date +%Y%m%d_%H%M%S)"
    print_info "Backing up existing files to: $BACKUP_DIR"
    cp -r "$TARGET_DIR" "$BACKUP_DIR"
    print_status "Backup created"
fi

# Clear target directory
print_info "Clearing target directory..."
rm -rf "${TARGET_DIR:?}"/*

# Copy build output to target directory
print_info "Copying Flutter web build to target directory..."
if cp -r "$SOURCE_DIR"/* "$TARGET_DIR/"; then
    print_status "Files copied successfully"
else
    print_error "Failed to copy files"
    exit 1
fi

# Verify copy was successful
file_count=$(find "$TARGET_DIR" -type f | wc -l)
if [ "$file_count" -gt 0 ]; then
    print_status "Build deployment completed successfully"
    print_info "Copied $file_count files to $TARGET_DIR"
else
    print_error "No files were copied - something went wrong"
    exit 1
fi

# List key files to verify
echo ""
print_info "Key files in target directory:"
if [ -f "$TARGET_DIR/index.html" ]; then
    print_status "index.html"
else
    print_warning "index.html not found"
fi

if [ -f "$TARGET_DIR/main.dart.js" ]; then
    print_status "main.dart.js"
else
    print_warning "main.dart.js not found"
fi

if [ -f "$TARGET_DIR/flutter.js" ]; then
    print_status "flutter.js"
else
    print_warning "flutter.js not found"
fi

if [ -d "$TARGET_DIR/assets" ]; then
    asset_count=$(find "$TARGET_DIR/assets" -type f | wc -l)
    print_status "assets/ directory ($asset_count files)"
else
    print_warning "assets/ directory not found"
fi

# Show directory size
target_size=$(du -sh "$TARGET_DIR" | cut -f1)
print_info "Total size of deployed files: $target_size"

echo ""
print_status "Flutter web UI build and deployment completed!"
print_info "The web application is now ready to serve from: $TARGET_DIR"
print_info "You can now run the SCVis Go application to serve the updated UI"

echo ""
print_info "Next steps:"
echo "  1. Start the SCVis Go application: go run ./cmd/web/"
echo "  2. Navigate to the visual interface in your browser"
echo "  3. Verify the updated UI is working correctly"
