#!/bin/bash

# Quick Build Script for SCVis (simplified version)
# Builds Flutter web app and Go backend

set -e

echo "🚀 Quick SCVis Build..."

# Build Flutter
echo "📱 Building Flutter web app..."
cd ../scvis-flutter && ./build_web.sh

# Copy to Go project
echo "📁 Copying web files..."
cd - > /dev/null
cp -r ../scvis-flutter/scvis_flutter/build/web/* visual/scvis/

# Build Go app
echo "🏗️  Building Go application..."
go build -o scvisual ./cmd/web

echo "✅ Build complete! Binary: ./scvisual"
