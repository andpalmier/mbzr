package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andpalmier/mbzr/api"
)

// Global flags
var verbose bool

// printRootHelp displays the help message for the root command
func printRootHelp() {
	fmt.Println("mbzr - MalwareBazaar CLI Client")
	fmt.Println("  A command-line tool for interacting with the MalwareBazaar API")
	fmt.Println("  Built by @andpalmier")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  mbzr [command] [flags]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  query              Query MalwareBazaar for malware samples")
	fmt.Println("  download           Download a malware sample by SHA256 hash")
	fmt.Println("  upload             Upload a file or directory to MalwareBazaar")
	fmt.Println("  comment            Add a comment to a malware sample")
	fmt.Println("  update             Update metadata of a malware sample")
	fmt.Println("  recent_detections  Get recent malware detections")
	fmt.Println("  cscb               Query the Code Signing Certificate Blocklist")
	fmt.Println("  version            Show version information")
	fmt.Println()
	fmt.Println("Global Flags:")
	fmt.Println("  -v, --verbose      Enable verbose output")
	fmt.Println("  -V, --version      Show version information")
	fmt.Println("  -h, --help         Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Query samples by tag")
	fmt.Println("  mbzr query -tag Emotet -limit 10")
	fmt.Println()
	fmt.Println("  # Download a sample")
	fmt.Println("  mbzr download -hash ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652")
	fmt.Println()
	fmt.Println("  # Upload a file with tags")
	fmt.Println("  mbzr upload -file malware.exe -tags trojan,banker")
	fmt.Println()
	fmt.Println("  # Get recent detections from last 24 hours")
	fmt.Println("  mbzr recent_detections -hours 24")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  MALWAREBAZAAR_API_KEY  Your MalwareBazaar API key (required for most commands)")
	fmt.Println()
	fmt.Println("For more information about a command:")
	fmt.Println("  mbzr [command] --help")
}

// getAPIClient creates and returns an API client with the API key from environment
// Returns an error if the API key is not set
func getAPIClient() (*api.Client, error) {
	apiKey := os.Getenv("MALWAREBAZAAR_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MALWAREBAZAAR_API_KEY environment variable is required")
	}

	if verbose {
		printVerbose("Creating API client")
	}

	return api.NewClient(apiKey), nil
}

// getContext returns a context with timeout and its cancel function
// Default timeout is 30 seconds for API requests
// Callers must defer cancel() to avoid context leaks
func getContext() (context.Context, context.CancelFunc) {
	timeout := 30 * time.Second

	if verbose {
		printVerbose(fmt.Sprintf("Setting request timeout to %v", timeout))
	}

	return context.WithTimeout(context.Background(), timeout)
}

// printUsageHeader prints a standard usage header for commands
func printUsageHeader(command, description string) {
	fmt.Printf("Usage:\n  mbzr %s [flags]\n", command)
	fmt.Println(description)
}

// printError prints an error message
func printError(message string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", message)
}

// printDetailedError prints a detailed error message with context
func printDetailedError(err error, context string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	if context != "" {
		fmt.Fprintf(os.Stderr, "Context: %s\n", context)
	}

	// Suggest solutions for common errors
	errStr := err.Error()
	suggestions := map[string]string{
		"Unauthorized":       "Set MALWAREBAZAAR_API_KEY environment variable\n          export MALWAREBAZAAR_API_KEY=your_key_here",
		"API key":            "Set MALWAREBAZAAR_API_KEY environment variable\n          export MALWAREBAZAAR_API_KEY=your_key_here",
		"timeout":            "The request timed out. Try again or check your network connection",
		"deadline exceeded":  "The request timed out. Try again or check your network connection",
		"connection refused": "Cannot reach API. Check your internet connection",
	}

	for keyword, solution := range suggestions {
		if contains(errStr, keyword) {
			fmt.Fprintf(os.Stderr, "Solution: %s\n", solution)
			break
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Full error: %+v\n", err)
	}
}

// printVerbose prints a verbose message
func printVerbose(message string) {
	fmt.Printf("[VERBOSE] %s\n", message)
}

// printSuccess prints a success message
func printSuccess(message string) {
	fmt.Println(message)
}

// printJSON prints data as indented JSON
func printJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		// Fallback to standard JSON if coloring fails
		fmt.Printf("%+v\n", data)
		return
	}
	fmt.Println(string(b))
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// SetVerbose sets the verbose flag (for testing)
func SetVerbose(v bool) {
	verbose = v
	InitLogger(v)
}

// IsVerbose returns the verbose flag value
func IsVerbose() bool {
	return verbose
}
