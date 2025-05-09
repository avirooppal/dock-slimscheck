package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/avirooppal/dock-slimscheck/checks"
	"github.com/avirooppal/dock-slimscheck/parser"
	"github.com/avirooppal/dock-slimscheck/security"
)

// Version information
const (
	Version = "1.0.0"
)

func main() {
	// Define command line flags
	securityFlag := flag.Bool("security", false, "Enable additional security checks")
	versionFlag := flag.Bool("version", false, "Display version information")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("dock-slimcheck version %s\n", Version)
		return
	}

	// Check if Dockerfile path is provided
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: No Dockerfile specified")
		fmt.Println("Usage: dock-slimcheck ./Dockerfile [--security]")
		os.Exit(1)
	}

	dockerfilePath := args[0]

	// Check if the Dockerfile exists
	_, err := os.Stat(dockerfilePath)
	if os.IsNotExist(err) {
		fmt.Printf("Error: Dockerfile not found at %s\n", dockerfilePath)
		os.Exit(1)
	}

	// Parse the Dockerfile
	fmt.Printf("[INFO] Checking Dockerfile: %s\n\n", dockerfilePath)
	dockerfile, err := parser.ParseDockerfile(dockerfilePath)
	if err != nil {
		fmt.Printf("Error parsing Dockerfile: %s\n", err)
		os.Exit(1)
	}

	// Determine the directory of the Dockerfile for contextual checks
	dockerfileDir := filepath.Dir(dockerfilePath)

	// Run checks and collect issues
	var issues []checks.Issue

	// Base image checks
	baseImageIssues := checks.CheckBaseImage(dockerfile)
	issues = append(issues, baseImageIssues...)

	// Best practices checks
	practiceIssues := checks.CheckBestPractices(dockerfile, dockerfileDir)
	issues = append(issues, practiceIssues...)

	// Layer size checks (requires Docker to be installed)
	if checks.IsDockerAvailable() {
		sizeIssues := checks.CheckLayerSizes(dockerfile, dockerfileDir)
		issues = append(issues, sizeIssues...)
	}

	// Security checks (if enabled)
	if *securityFlag {
		securityIssues := security.RunSecurityChecks(dockerfile)
		issues = append(issues, securityIssues...)
	}

	// Print issues
	printIssues(issues)

	// Return non-zero exit code if issues were found
	if len(issues) > 0 {
		os.Exit(2)
	}
}

func printIssues(issues []checks.Issue) {
	// Set up colors
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	for _, issue := range issues {
		// Print issue header with a newline before each issue
		fmt.Println()
		prefix := yellow("[!]")
		if issue.Type == checks.SecurityIssue {
			prefix = red("[SECURITY]")
		} else if issue.Type == checks.InfoIssue {
			prefix = green("[+]")
		}

		fmt.Printf("%s %s\n", prefix, issue.Message)
		
		// Print severity and impact if they exist
		if issue.Severity != "" {
			fmt.Printf("  %s Severity: %s\n", blue("→"), issue.Severity)
		}
		if issue.Impact != "" {
			fmt.Printf("  %s Impact: %s\n", blue("→"), issue.Impact)
		}
		
		// Print fix suggestion with proper indentation
		if issue.Fix != "" {
			fmt.Printf("  %s Fix:\n", blue("→"))
			// Split fix into lines and indent each line
			fixLines := strings.Split(issue.Fix, "\n")
			for _, line := range fixLines {
				fmt.Printf("    %s\n", cyan(line))
			}
		}
		
		// Print references
		if len(issue.References) > 0 {
			fmt.Printf("  %s References:\n", blue("→"))
			for _, ref := range issue.References {
				fmt.Printf("    - %s\n", ref)
			}
		}
	}

	// Print summary with a newline before it
	fmt.Println()
	fmt.Printf("[%s] Check complete — %d issues found\n", green("✓"), len(issues))
}