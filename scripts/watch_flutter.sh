#!/bin/bash

# Flutter Development Watch Script
# Watches for changes in Flutter project and automatically rebuilds

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

FLUTTER_DIR="../scvis_flutter"
TARGET_DIR="./visual/scvis"

echo "========================================"
echo "    Flutter Development Watch Mode"
echo "========================================"
echo ""

# Check if inotify-tools is installed
if ! command -v inotifywait &> /dev/null; then
    print_warning "inotify-tools not found. Installing for file watching..."
    echo "You may need to run: sudo apt-get install inotify-tools"
    echo "For now, falling back to polling mode..."
    POLLING_MODE=true
else
    print_status "inotify-tools found"
    POLLING_MODE=false
fi

# Check if Flutter project exists
if [ ! -d "$FLUTTER_DIR" ]; then
    echo "Error: Flutter project not found at $FLUTTER_DIR"
    exit 1
fi

print_status "Flutter project found at: $FLUTTER_DIR"
print_info "Target directory: $TARGET_DIR"
print_info "Press Ctrl+C to stop watching"
echo ""

# Initial build
print_info "Performing initial build..."
./build_ui.sh

# Function to rebuild
rebuild() {
    print_info "Changes detected, rebuilding..."
    if ./build_ui.sh; then
        print_status "Rebuild completed at $(date '+%H:%M:%S')"
    else
        print_warning "Rebuild failed at $(date '+%H:%M:%S')"
    fi
    echo ""
}

if [ "$POLLING_MODE" = true ]; then
    # Polling mode - check for changes every 2 seconds
    print_info "Using polling mode (checking every 2 seconds)..."
    last_change=$(find "$FLUTTER_DIR" -name "*.dart" -o -name "*.yaml" -o -name "*.json" | xargs stat -c %Y | sort -n | tail -1)
    
    while true; do
        sleep 2
        current_change=$(find "$FLUTTER_DIR" -name "*.dart" -o -name "*.yaml" -o -name "*.json" | xargs stat -c %Y | sort -n | tail -1)
        if [ "$current_change" != "$last_change" ]; then
            rebuild
            last_change=$current_change
        fi
    done
else
    # inotify mode - real-time file watching
    print_info "Watching for changes in: $FLUTTER_DIR"
    print_info "Monitoring: *.dart, *.yaml, *.json files"
    
    while true; do
        inotifywait -r -e modify,create,delete \
            --include '\.(dart|yaml|json)$' \
            "$FLUTTER_DIR" 2>/dev/null
        
        # Debounce - wait a moment for multiple rapid changes
        sleep 1
        rebuild
    done
fi
