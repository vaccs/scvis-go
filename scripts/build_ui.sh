#!/bin/bash

# Quick Flutter Build Script
# Simple script to build and deploy Flutter web UI

set -e

echo "Building Flutter web UI..."

# Check if Flutter project exists
if [ ! -d "../scvis-flutter" ]; then
    echo "Error: Flutter project not found at ../scvis-flutter"
    exit 1
fi

# Build Flutter web app
cd ../scvis-flutter
if [ -f "build_web.sh" ]; then
    echo "Running build_web.sh..."
    ./build_web.sh
else
    echo "build_web.sh not found, running flutter build web directly..."
    flutter build web
fi

# Copy to target directory
cd - > /dev/null
echo "Copying build output..."
rm -rf visual/scvis/*
cp -r ../scvis-flutter/scvs_flutter/build/web/* visual/scvis/

echo "✓ Flutter web UI built and deployed successfully!"
echo "Files copied to: visual/scvis/"
