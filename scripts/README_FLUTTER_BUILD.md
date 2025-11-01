# Flutter UI Build System

This directory contains scripts to build and deploy the Flutter web UI for SCVis.

## Overview

The SCVis application consists of two main components:
- **Go Backend**: The main SCVis application (`cmd/web/`)
- **Flutter Frontend**: The web UI located at `../scvis_flutter`

The build system automates the process of building the Flutter web application and copying the build output to the correct location for serving by the Go backend.

## Available Scripts

### 1. `build_ui.sh` (Simple Build)

**Quick and simple build script for immediate use.**

```bash
./build_ui.sh
```

**What it does:**
- Checks if `../scvis_flutter` exists
- Runs `build_web.sh` in the Flutter project (or `flutter build web` if build script not found)
- Copies build output from `../scvis_flutter/build/web/` to `./visual/scvis/`
- Provides simple success/error feedback

**Use this when:**
- You want a quick rebuild after making Flutter changes
- You're doing one-time builds
- You want minimal output and fast execution

### 2. `scripts/build_flutter_ui.sh` (Comprehensive Build)

**Full-featured build script with detailed output and error handling.**

```bash
./scripts/build_flutter_ui.sh
```

**What it does:**
- Comprehensive pre-build checks and validation
- Backs up existing UI files before replacing
- Detailed progress reporting with colored output
- Verifies build success and file integrity
- Provides deployment statistics and next steps
- Enhanced error handling and troubleshooting info

**Use this when:**
- You want detailed feedback about the build process
- You're troubleshooting build issues
- You want automatic backups of previous builds
- You're doing production deployments

### 3. `scripts/watch_flutter.sh` (Development Watch Mode)

**Automated development workflow that watches for changes and rebuilds automatically.**

```bash
./scripts/watch_flutter.sh
```

**What it does:**
- Performs initial build
- Monitors Flutter project for file changes (`*.dart`, `*.yaml`, `*.json`)
- Automatically rebuilds and deploys when changes are detected
- Uses `inotify` for efficient file watching (falls back to polling if not available)
- Provides timestamped build notifications
- Runs until manually stopped with Ctrl+C

**Use this when:**
- You're actively developing the Flutter UI
- You want automatic rebuilds during development
- You want to see changes immediately in the browser

## Project Structure

```
scvis-go/                    # Main Go project
├── visual/scvis/           # Deployed Flutter web files (target)
├── build_ui.sh             # Simple build script
└── scripts/
    ├── build_flutter_ui.sh  # Comprehensive build script
    └── watch_flutter.sh     # Development watch script

../scvis_flutter/           # Flutter project (source)
├── build_web.sh           # Flutter build script
└── build/web/             # Flutter build output
```

## Requirements

### For All Scripts:
- Flutter project must be located at `../scvis_flutter`
- The Flutter project should have a `build_web.sh` script (or Flutter CLI available)

### For Watch Mode:
- **Recommended**: `inotify-tools` package for efficient file watching
  ```bash
  sudo apt-get install inotify-tools
  ```
- **Fallback**: Polling mode if inotify-tools not available

## Usage Examples

### Basic Development Workflow

1. **Initial Build**:
   ```bash
   ./build_ui.sh
   ```

2. **Start Go Backend**:
   ```bash
   go run ./cmd/web/
   ```

3. **For Active Development** (in another terminal):
   ```bash
   ./scripts/watch_flutter.sh
   ```

### Production Deployment

1. **Full Build with Verification**:
   ```bash
   ./scripts/build_flutter_ui.sh
   ```

2. **Build Go Application**:
   ```bash
   go build ./cmd/web/
   ```

3. **Deploy Both Components**

### Troubleshooting Build Issues

1. **Use Comprehensive Build Script**:
   ```bash
   ./scripts/build_flutter_ui.sh
   ```
   This provides detailed error messages and diagnostics.

2. **Check Flutter Project**:
   ```bash
   cd ../scvis_flutter
   flutter doctor
   flutter clean
   flutter pub get
   ```

3. **Verify Paths**:
   - Ensure `../scvis_flutter` exists and contains the Flutter project
   - Check that `../scvis_flutter/build_web.sh` exists and is executable
   - Verify `./visual/scvis/` directory is writable

## Build Output

After a successful build, the `visual/scvis/` directory will contain:

```
visual/scvis/
├── index.html              # Main HTML file
├── main.dart.js            # Compiled Dart code
├── flutter.js              # Flutter web framework
├── flutter_service_worker.js
├── assets/                 # Application assets
├── canvaskit/             # Flutter rendering engine
└── icons/                 # Application icons
```

## Integration with Go Backend

The Go backend serves the Flutter UI through:
- **Static file serving**: `visual/scvis/` directory
- **Route**: `/scvis/visual/` endpoint
- **Integration**: The `projectVisual` handler redirects to the Flutter app

## Development Tips

### Efficient Development Workflow

1. **Use Watch Mode**: Start `./scripts/watch_flutter.sh` once and leave it running
2. **Hot Reload**: The watch script automatically rebuilds on changes
3. **Browser Refresh**: Refresh your browser to see changes after rebuild
4. **Terminal Management**: Use separate terminals for watch mode and Go backend

### Build Optimization

1. **Clean Builds**: Occasionally run `flutter clean` in the Flutter project
2. **Dependency Updates**: Keep Flutter dependencies up to date
3. **Build Caching**: Flutter builds are incremental, so subsequent builds are faster

### Debugging Build Issues

1. **Check Flutter Environment**:
   ```bash
   cd ../scvis_flutter
   flutter doctor
   ```

2. **Manual Build Test**:
   ```bash
   cd ../scvis_flutter
   flutter build web --verbose
   ```

3. **Check File Permissions**:
   ```bash
   ls -la ../scvis_flutter/build_web.sh
   ls -la ./visual/scvis/
   ```

## Customization

### Custom Build Arguments

Modify the scripts to pass custom arguments to Flutter build:

```bash
# In build_ui.sh or build_flutter_ui.sh
flutter build web --release --dart-define=ENVIRONMENT=production
```

### Different Flutter Project Location

If your Flutter project is in a different location, update the `FLUTTER_PROJECT_DIR` variable in the scripts:

```bash
FLUTTER_PROJECT_DIR="/path/to/your/flutter/project"
```

### Custom Target Directory

To deploy to a different directory, update the `TARGET_DIR` variable:

```bash
TARGET_DIR="./custom/ui/path"
```

## Troubleshooting

### Common Issues

1. **"Flutter project not found"**
   - Ensure `../scvis_flutter` directory exists
   - Check the relative path from your current directory

2. **"build_web.sh not found"**
   - Create a `build_web.sh` script in the Flutter project
   - Or the script will fall back to `flutter build web`

3. **"Permission denied"**
   - Make scripts executable: `chmod +x script_name.sh`
   - Check write permissions for target directory

4. **"No files copied"**
   - Verify Flutter build completed successfully
   - Check if `../scvis_flutter/build/web/` contains files

### Getting Help

- Check script output for specific error messages
- Use the comprehensive build script for detailed diagnostics
- Verify Flutter environment with `flutter doctor`
- Ensure all prerequisites are installed
