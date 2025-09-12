package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"scvis.vaccsedu.org/internal/vaccs"
)

// Configuration structure for command line arguments
type AnalysisConfig struct {
	Files           []string // C files to analyze
	VaccsArgs       []string // Arguments for the VACCS analyzer
	ExecutableArgs  []string // Arguments for the C program execution
	ExecutableName  string   // Name of the executable to create
	OutputFile      string   // Output file for analysis results
	OutputFormat    string   // Output format: json, pretty, summary
	Verbose         bool     // Verbose logging
	CheckConnection bool     // Only check if VACCS server is running
}

func main() {
	config := parseFlags()

	// Setup logging
	logLevel := slog.LevelInfo
	if config.Verbose {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// Create VACCS communication instance
	vc := vaccs.NewVaccsComm(logger)

	// Check connection if requested
	if config.CheckConnection {
		if vc.IsPortOpen() {
			fmt.Println("✓ VACCS server is running on localhost:3580")
			os.Exit(0)
		} else {
			fmt.Println("✗ VACCS server is not running on localhost:3580")
			os.Exit(1)
		}
	}

	// Validate input files
	if len(config.Files) == 0 {
		log.Fatal("No C files specified for analysis")
	}

	// Read and validate C files
	filenames := make([]string, len(config.Files))
	sources := make([]string, len(config.Files))

	for i, filename := range config.Files {
		source, err := readCFile(filename)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", filename, err)
		}
		filenames[i] = filepath.Base(filename)
		sources[i] = source
		logger.Info("Loaded C file", "filename", filenames[i], "size", len(source))
	}

	// Set default executable name if not provided
	if config.ExecutableName == "" {
		config.ExecutableName = "analysis_program"
	}

	logger.Info("Starting analysis",
		"files", len(filenames),
		"vaccs_args", len(config.VaccsArgs),
		"exec_args", len(config.ExecutableArgs),
		"executable", config.ExecutableName)

	// Connect to VACCS server
	err := vc.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to VACCS server: %v", err)
	}
	defer vc.Close()

	logger.Info("Connected to VACCS server")

	// Perform analysis
	result, err := vc.PerformAnalysis(filenames, sources, config.VaccsArgs, config.ExecutableName, config.ExecutableArgs)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	logger.Info("Analysis completed", "result_size", len(result))

	// Output results
	err = outputResults(result, config)
	if err != nil {
		log.Fatalf("Failed to output results: %v", err)
	}

	fmt.Println("✓ Analysis completed successfully")
}

func parseFlags() *AnalysisConfig {
	config := &AnalysisConfig{}

	// Define flags
	var vaccsArgs string
	var execArgs string

	flag.StringVar(&vaccsArgs, "vaccs-args", "", "Comma-separated list of arguments for VACCS analyzer (e.g., \"-Wall,-O2\")")
	flag.StringVar(&execArgs, "exec-args", "", "Comma-separated list of arguments for C program execution (e.g., \"arg1,arg2,arg3\")")
	flag.StringVar(&config.ExecutableName, "executable", "", "Name of the executable to create (default: analysis_program)")
	flag.StringVar(&config.OutputFile, "output", "", "Output file for analysis results (default: stdout)")
	flag.StringVar(&config.OutputFormat, "format", "pretty", "Output format: json, pretty, summary")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&config.CheckConnection, "check", false, "Only check if VACCS server is running")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file1.c> [file2.c] ...\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Analyze C files using the VACCS analyzer.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s program.c\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -vaccs-args=\"-Wall,-O2\" -exec-args=\"10,20\" program.c\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -format=json -output=results.json program.c library.c\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -check\n", os.Args[0])
	}

	flag.Parse()

	// Parse remaining arguments as C files
	config.Files = flag.Args()

	// Parse comma-separated arguments
	if vaccsArgs != "" {
		config.VaccsArgs = strings.Split(vaccsArgs, ",")
		// Trim whitespace from each argument
		for i, arg := range config.VaccsArgs {
			config.VaccsArgs[i] = strings.TrimSpace(arg)
		}
	}

	if execArgs != "" {
		config.ExecutableArgs = strings.Split(execArgs, ",")
		// Trim whitespace from each argument
		for i, arg := range config.ExecutableArgs {
			config.ExecutableArgs[i] = strings.TrimSpace(arg)
		}
	}

	// Validate output format
	switch config.OutputFormat {
	case "json", "pretty", "summary":
		// Valid formats
	default:
		fmt.Fprintf(os.Stderr, "Invalid output format: %s\n", config.OutputFormat)
		fmt.Fprintf(os.Stderr, "Valid formats: json, pretty, summary\n")
		os.Exit(1)
	}

	return config
}

func readCFile(filename string) (string, error) {
	// Check if file exists and is readable
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filename)
	}

	// Read file contents
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Validate it looks like a C file
	contentStr := string(content)
	if !strings.Contains(contentStr, "#include") && !strings.Contains(contentStr, "int main") && !strings.Contains(contentStr, "void") {
		fmt.Fprintf(os.Stderr, "Warning: %s does not appear to be a C source file\n", filename)
	}

	return contentStr, nil
}

func outputResults(result []byte, config *AnalysisConfig) error {
	var output io.Writer = os.Stdout

	// Open output file if specified
	if config.OutputFile != "" {
		file, err := os.Create(config.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer file.Close()
		output = file
	}

	switch config.OutputFormat {
	case "json":
		// Output raw JSON
		_, err := output.Write(result)
		return err

	case "pretty":
		// Pretty-print JSON
		var jsonData interface{}
		if err := json.Unmarshal(result, &jsonData); err != nil {
			return fmt.Errorf("failed to parse JSON result: %v", err)
		}

		prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format JSON: %v", err)
		}

		_, err = output.Write(prettyJSON)
		return err

	case "summary":
		// Output a summary of the analysis
		return outputSummary(result, output)

	default:
		return fmt.Errorf("unknown output format: %s", config.OutputFormat)
	}
}

func outputSummary(result []byte, output io.Writer) error {
	// Parse JSON to extract summary information
	var analysis map[string]interface{}
	if err := json.Unmarshal(result, &analysis); err != nil {
		return fmt.Errorf("failed to parse analysis result: %v", err)
	}

	fmt.Fprintf(output, "VACCS Analysis Summary\n")
	fmt.Fprintf(output, "=====================\n\n")

	// Try to extract common fields from the analysis
	if timestamp, ok := analysis["timestamp"]; ok {
		fmt.Fprintf(output, "Analysis Time: %v\n", timestamp)
	}

	if status, ok := analysis["status"]; ok {
		fmt.Fprintf(output, "Status: %v\n", status)
	}

	if errors, ok := analysis["errors"]; ok {
		if errList, ok := errors.([]interface{}); ok {
			fmt.Fprintf(output, "Errors Found: %d\n", len(errList))
			for i, err := range errList {
				fmt.Fprintf(output, "  %d. %v\n", i+1, err)
			}
		}
	}

	if warnings, ok := analysis["warnings"]; ok {
		if warnList, ok := warnings.([]interface{}); ok {
			fmt.Fprintf(output, "Warnings Found: %d\n", len(warnList))
			for i, warn := range warnList {
				fmt.Fprintf(output, "  %d. %v\n", i+1, warn)
			}
		}
	}

	// Count total issues
	totalIssues := 0
	if errors, ok := analysis["errors"].([]interface{}); ok {
		totalIssues += len(errors)
	}
	if warnings, ok := analysis["warnings"].([]interface{}); ok {
		totalIssues += len(warnings)
	}

	fmt.Fprintf(output, "\nTotal Issues: %d\n", totalIssues)

	// Show raw JSON size
	fmt.Fprintf(output, "Raw Analysis Size: %d bytes\n", len(result))

	if totalIssues == 0 {
		fmt.Fprintf(output, "\n✓ No issues found in the analyzed code.\n")
	} else {
		fmt.Fprintf(output, "\n⚠ Issues found. Review the detailed analysis for more information.\n")
	}

	return nil
}
