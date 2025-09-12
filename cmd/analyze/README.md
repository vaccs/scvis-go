# VACCS Analysis Tool

A command-line tool for analyzing C files using the VACCS analyzer.

## Usage

```bash
go run ./cmd/analyze [options] <file1.c> [file2.c] ...
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `-vaccs-args` | Comma-separated list of arguments for VACCS analyzer | "" |
| `-exec-args` | Comma-separated list of arguments for C program execution | "" |
| `-executable` | Name of the executable to create | "analysis_program" |
| `-output` | Output file for analysis results | stdout |
| `-format` | Output format: json, pretty, summary | "pretty" |
| `-verbose` | Enable verbose logging | false |
| `-check` | Only check if VACCS server is running | false |

## Examples

### Basic Analysis
```bash
# Analyze a single C file
go run ./cmd/analyze program.c

# Analyze multiple C files
go run ./cmd/analyze main.c utils.c helper.c
```

### With VACCS Arguments
```bash
# Pass compiler flags to VACCS
go run ./cmd/analyze -vaccs-args="-Wall,-O2,-std=c99" program.c

# Pass analysis-specific arguments
go run ./cmd/analyze -vaccs-args="-analyzer-checker=security" program.c
```

### With Execution Arguments
```bash
# Provide runtime arguments for the C program
go run ./cmd/analyze -exec-args="10,20,30" calculator.c

# Provide string arguments
go run ./cmd/analyze -exec-args="input.txt,output.txt" fileprocessor.c
```

### Output Formats

#### Pretty JSON (Default)
```bash
go run ./cmd/analyze -format=pretty program.c
```

#### Raw JSON
```bash
go run ./cmd/analyze -format=json program.c
```

#### Summary
```bash
go run ./cmd/analyze -format=summary program.c
```

### Save to File
```bash
# Save results to a file
go run ./cmd/analyze -output=analysis_results.json -format=json program.c

# Save summary to file
go run ./cmd/analyze -output=summary.txt -format=summary program.c
```

### Verbose Logging
```bash
# Enable detailed logging
go run ./cmd/analyze -verbose program.c
```

### Check Server Connection
```bash
# Check if VACCS server is running
go run ./cmd/analyze -check
```

## Complete Examples

### Simple Security Analysis
```bash
go run ./cmd/analyze \
  -vaccs-args="-analyzer-checker=security.insecureAPI.UncheckedReturn" \
  -format=summary \
  -output=security_report.txt \
  network_server.c
```

### Multi-file Project Analysis
```bash
go run ./cmd/analyze \
  -vaccs-args="-Wall,-Wextra,-O2" \
  -exec-args="config.ini,debug" \
  -executable="myprogram" \
  -format=pretty \
  -output=full_analysis.json \
  main.c utils.c network.c database.c
```

### Quick Check with Verbose Output
```bash
go run ./cmd/analyze -verbose -format=summary small_program.c
```

## Building the Tool

To build a standalone executable:

```bash
# Build for current platform
go build -o vaccs-analyze ./cmd/analyze

# Use the built executable
./vaccs-analyze program.c
```

### Cross-platform builds
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o vaccs-analyze-linux ./cmd/analyze

# Windows
GOOS=windows GOARCH=amd64 go build -o vaccs-analyze.exe ./cmd/analyze

# macOS
GOOS=darwin GOARCH=amd64 go build -o vaccs-analyze-mac ./cmd/analyze
```

## Prerequisites

1. **VACCS Server**: The VACCS analysis server must be running on `localhost:3580`
2. **Go Environment**: Go 1.19 or later
3. **C Files**: Valid C source files for analysis

## Output Formats

### JSON Format
Raw JSON output from the VACCS analyzer, suitable for programmatic processing.

### Pretty Format
Formatted JSON output that's human-readable.

### Summary Format
A condensed summary showing:
- Analysis timestamp and status
- Number of errors and warnings found
- Brief description of each issue
- Overall assessment

## Error Handling

The tool will exit with appropriate error codes:
- `0`: Success
- `1`: Connection failure, file not found, or analysis errors

Common issues:
- **"Failed to connect to VACCS server"**: Ensure the VACCS server is running on port 3580
- **"File does not exist"**: Check file paths are correct
- **"Analysis failed"**: Check the C code for syntax errors or unsupported constructs

## Integration with SCVis

This tool can be integrated into the SCVis web application workflow or used standalone for command-line analysis. The JSON output is compatible with the SCVis project analysis format.
