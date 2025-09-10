# SCVis Build Scripts

This directory contains build scripts for the SCVis application, which combines a Flutter web frontend with a Go backend.

## Available Build Scripts

### 1. `build.sh` (Comprehensive Build)

The main build script that performs a complete build process with detailed output and verification.

**What it does:**
1. Runs `../scvis-flutter/build_web.sh` to build the Flutter web application
2. Copies the built Flutter files to `visual/scvis/` in this project
3. Builds the Go application (`scvisual` binary)
4. Performs verification and shows detailed build summary

**Usage:**
```bash
./build.sh
```

**Features:**
- Colored output for better readability
- Detailed progress reporting
- Error handling with proper exit codes
- Build verification
- File count and size reporting
- Comprehensive build summary

### 2. `quick-build.sh` (Fast Build)

A simplified version for quick builds during development.

**Usage:**
```bash
./quick-build.sh
```

**Features:**
- Minimal output
- Fast execution
- Same functionality as the comprehensive build
- Ideal for development iterations

## Prerequisites

### Flutter Setup
- Flutter SDK installed and in PATH
- Flutter project located at `../scvis-flutter/`
- The `build_web.sh` script must exist in the Flutter project

### Go Setup
- Go installed and configured
- Valid `go.mod` file in this directory
- Go source code in `cmd/web/`

## Build Output

### Flutter Web Files
- **Location:** `visual/scvis/`
- **Key files:** `index.html`, `main.dart.js`, `assets/`, etc.
- **Purpose:** Static web files served by the Go application

### Go Binary
- **Location:** `./scvisual`
- **Purpose:** Main application executable
- **Size:** ~12MB (includes embedded static files)

## Directory Structure

```
scvis-go/                    # This project
├── build.sh                 # Main build script
├── quick-build.sh           # Quick build script
├── visual/scvis/            # Flutter web files (build output)
├── cmd/web/                 # Go source code
└── scvisual                 # Go binary (after build)

../scvis-flutter/            # Flutter project
├── build_web.sh            # Flutter build script
└── scvis_flutter/          # Flutter source code
    └── build/web/          # Flutter build output
```

## Usage Examples

### Development Workflow
```bash
# Quick build during development
./quick-build.sh

# Test the application
./scvisual

# Or deploy to production
sudo ./deploy.sh
```

### Production Build
```bash
# Comprehensive build with verification
./build.sh

# Deploy the built application
sudo ./deploy.sh

# Check deployment status
sudo systemctl status scvisual
```

### Build Only Go (Skip Flutter)
```bash
# If Flutter files are already up to date
go build -o scvisual ./cmd/web
```

### Build Only Flutter (Skip Go)
```bash
# From the Flutter directory
cd ../scvis-flutter && ./build_web.sh

# Then manually copy if needed
cp -r ../scvis-flutter/scvis_flutter/build/web/* visual/scvis/
```

## Troubleshooting

### Flutter Build Issues
- Ensure Flutter is installed: `flutter --version`
- Check Flutter project: `cd ../scvis-flutter && flutter doctor`
- Clean Flutter cache: `cd ../scvis-flutter/scvis_flutter && flutter clean`

### Go Build Issues
- Verify Go installation: `go version`
- Check dependencies: `go mod tidy`
- Clean Go cache: `go clean -cache`

### File Copy Issues
- Check permissions on `visual/scvis/` directory
- Ensure Flutter build completed successfully
- Verify source files exist in `../scvis-flutter/scvis_flutter/build/web/`

## Integration with Deployment

After building, you can deploy using:
```bash
sudo ./deploy.sh
```

This will:
- Install the binary to `/opt/scvisual/`
- Copy static files including the built Flutter web app
- Configure and start the systemd service
- Set up proper permissions and security

## Build Information

- **Flutter Build Mode:** Release (optimized for production)
- **Go Build Target:** Native binary for current platform
- **Tree Shaking:** Enabled for Flutter (reduces bundle size)
- **Static Files:** Embedded in Go binary via `embed` package
