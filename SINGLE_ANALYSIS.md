# Single Analysis Feature

## Overview
The Single Analysis feature allows users to quickly analyze C files without requiring login or database storage. This is ideal for one-off analyses or quick testing.

## Access
- **URL**: `/scvis/single`
- **Authentication**: Not required
- **Database**: Not used

## Features

### 1. File Upload
- **C Files**: Required - Upload one or more .c files for analysis
- **Input Files**: Optional - Upload supporting files of any type (test data, config files, etc.)

### 2. Analysis Options
- **Title**: Optional name for your analysis session
- **Parameters**: Optional command-line parameters for the analyzer (e.g., `-cgo`, `-esp`, `-scc`)
- **Stdin File**: Select which input file to use as stdin for the analysis

### 3. How It Works

1. **Upload Files**
   - Select C files to analyze (required)
   - Optionally upload input files
   - Optionally specify analysis parameters

2. **Temporary Storage**
   - Files are stored in a temporary directory (`/tmp/scvis-single-*`)
   - Directory structure:
     ```
     /tmp/scvis-single-XXXXXX/
     ├── cfiles/
     │   ├── file1.c
     │   └── file2.c
     ├── inputfiles/
     │   ├── test_data.txt
     │   └── config.json
     └── analysis.json
     ```

3. **Analysis**
   - C files are sent to the analysis service
   - Analysis results are returned as JSON
   - Results are saved to the temporary directory

4. **Visualization**
   - Results are automatically forwarded to `/scvis/visual`
   - Visualization loads analysis data from the temporary directory
   - No database interaction required

### 4. Cleanup
- Temporary directories are cleaned up after visualization
- Directories may also be cleaned on server restart
- No persistent storage in the database

## Usage Example

### Simple Analysis
1. Navigate to `/scvis/single`
2. Upload your C file(s)
3. Click "Analyze"
4. View results in the visualization

### Analysis with Input File
1. Navigate to `/scvis/single`
2. Upload your C file(s)
3. Upload an input file (e.g., test data)
4. Select the input file as stdin
5. Add parameters if needed
6. Click "Analyze"
7. View results

## API Flow

```
User uploads files
    ↓
POST /scvis/single/analyze
    ↓
Create temp directory
    ↓
Save C files to temp/cfiles/
    ↓
Save input files to temp/inputfiles/
    ↓
Connect to analysis service (vaccs)
    ↓
Send files and parameters
    ↓
Receive analysis JSON
    ↓
Save to temp/analysis.json
    ↓
Redirect to /scvis/visual?tempdir=...&title=...
    ↓
Visualization loads from temp directory
```

## Technical Details

### Handler Functions
- `singleAnalyze()` - GET handler for the form
- `singleAnalyzePost()` - POST handler for file processing and analysis

### File Processing
- Maximum upload size: 32MB
- C files must have `.c` extension
- Input files can be any type
- Files are validated before processing

### Error Handling
- Connection failures to analysis service
- File read/write errors
- Analysis errors (compile, analyze, protocol, timeout, system)
- User-friendly error messages displayed

### Parameters
The following parameters are supported:
- `-cgo`: Enable C goto analysis
- `-esp`: Enable extended static path analysis
- `-scc`: Enable strongly connected components analysis
- And other analyzer-specific parameters

## Security Considerations

1. **No Authentication**: Anyone can access this feature
2. **Temporary Storage**: Files are stored temporarily and cleaned up
3. **No Database**: No user data or files are permanently stored
4. **File Size Limits**: 32MB maximum upload size
5. **CSRF Protection**: Enabled via noSurf middleware

## Differences from Project Analysis

| Feature | Project Analysis | Single Analysis |
|---------|-----------------|-----------------|
| Authentication | Required | Not required |
| Database Storage | Yes | No |
| File Persistence | Permanent | Temporary |
| Project Management | Full | None |
| User Tracking | Yes | No |
| Multiple Sessions | Yes | One-time |

## Future Enhancements

Possible improvements:
1. Add rate limiting to prevent abuse
2. Implement temporary file cleanup cron job
3. Add download option for analysis results
4. Support for sharing analysis links
5. Optional user account linking for saving results
