package checks

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/avirooppal/dock-slimscheck.git/parser"
)

// CheckBestPractices performs various best practice checks on the Dockerfile
func CheckBestPractices(dockerfile *parser.Dockerfile, contextDir string) []Issue {
	var issues []Issue

	// Check for .dockerignore when using COPY . .
	if dockerfile.HasWildcardCopy() {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "COPY . . used — consider using specific paths",
		})

		// Check if .dockerignore exists
		dockerignorePath := filepath.Join(contextDir, ".dockerignore")
		if _, err := os.Stat(dockerignorePath); os.IsNotExist(err) {
			issues = append(issues, Issue{
				Type:    WarningIssue,
				Message: "No `.dockerignore` found",
			})
		}
	}

	// Check for ADD vs COPY
	for _, instruction := range dockerfile.Instructions {
		if instruction.Command == "ADD" {
			issues = append(issues, Issue{
				Type:    WarningIssue,
				Message: "Using ADD instead of COPY — ADD adds unneeded complexity and risk",
			})
			break
		}
	}

	// Check if HEALTHCHECK is missing
	if !dockerfile.HasHealthcheck() {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "No HEALTHCHECK found",
		})
	}

	// Check if USER is missing
	if !dockerfile.HasUser() {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "USER not specified — running as root",
		})
	}

	// Check for cleanup after package installations
	issues = append(issues, checkPackageCleanup(dockerfile)...)

	// Check for latest tag in base image
	if strings.HasSuffix(dockerfile.BaseImage, ":latest") || !strings.Contains(dockerfile.BaseImage, ":") {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "Using ':latest' tag or no tag specified — this is non-reproducible",
		})
	}

	return issues
}

// checkPackageCleanup checks if apt/yum installations are properly cleaned up
func checkPackageCleanup(dockerfile *parser.Dockerfile) []Issue {
	var issues []Issue
	var hasAptGet, hasYum, hasApk bool
	var hasCleanup bool

	for _, instruction := range dockerfile.Instructions {
		if instruction.Command == "RUN" {
			// Check for package managers
			if strings.Contains(instruction.Arguments, "apt-get install") || 
			   strings.Contains(instruction.Arguments, "apt install") {
				hasAptGet = true
				
				// Check if cleanup exists in the same instruction
				if strings.Contains(instruction.Arguments, "apt-get clean") || 
				   strings.Contains(instruction.Arguments, "rm -rf /var/lib/apt/lists/*") {
					hasCleanup = true
				}
			}
			
			if strings.Contains(instruction.Arguments, "yum install") {
				hasYum = true
				
				// Check if cleanup exists in the same instruction
				if strings.Contains(instruction.Arguments, "yum clean all") {
					hasCleanup = true
				}
			}
			
			if strings.Contains(instruction.Arguments, "apk add") {
				hasApk = true
				
				// Check if cleanup exists in the same instruction
				if strings.Contains(instruction.Arguments, "--no-cache") || 
				   strings.Contains(instruction.Arguments, "rm -rf /var/cache/apk/*") {
					hasCleanup = true
				}
			}
		}
	}

	// If any package manager is used without cleanup
	if (hasAptGet || hasYum || hasApk) && !hasCleanup {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "Package installation without cleanup — adds unnecessary size",
		})
	}

	return issues
}