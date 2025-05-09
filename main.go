package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

	// Print summary
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("\n[%s] Check complete — %d issues found\n", green("✓"), len(issues))

	// Return non-zero exit code if issues were found
	if len(issues) > 0 {
		os.Exit(2)
	}
}

func printIssues(issues []checks.Issue) {
	// Set up colors
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for _, issue := range issues {
		prefix := yellow("[!]")
		if issue.Type == checks.SecurityIssue {
			prefix = red("[SECURITY]")
		} else if issue.Type == checks.InfoIssue {
			prefix = "[+]"
		}

		fmt.Printf("%s %s\n", prefix, issue.Message)
	}
}