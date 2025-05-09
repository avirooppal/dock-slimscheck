package checks

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/avirooppal/dock-slimscheck/parser"
)

// CheckBestPractices performs various best practice checks on the Dockerfile
func CheckBestPractices(dockerfile *parser.Dockerfile, contextDir string) []Issue {
	var issues []Issue

	// Check for .dockerignore when using COPY . .
	if dockerfile.HasWildcardCopy() {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "COPY . . used — consider using specific paths",
			Fix:     "Replace 'COPY . .' with specific paths, e.g., 'COPY package.json package-lock.json ./'",
			Severity: "medium",
			Impact:   "Large context size and potential inclusion of sensitive files",
			References: []string{
				"https://docs.docker.com/develop/dev-best-practices/#use-specific-paths",
			},
		})

		// Check if .dockerignore exists
		dockerignorePath := filepath.Join(contextDir, ".dockerignore")
		if _, err := os.Stat(dockerignorePath); os.IsNotExist(err) {
			issues = append(issues, Issue{
				Type:    WarningIssue,
				Message: "No `.dockerignore` found",
				Fix:     "Create a .dockerignore file with entries like:\nnode_modules\n.git\n*.md\n.env\n.DS_Store\ndist\nbuild\n*.log",
				Severity: "medium",
				Impact:   "Increased build context size and potential inclusion of sensitive files",
				References: []string{
					"https://docs.docker.com/develop/dev-best-practices/#use-dockerignore",
				},
			})
		}
	}

	// Check for ADD vs COPY
	for _, instruction := range dockerfile.Instructions {
		if instruction.Command == "ADD" {
			issues = append(issues, Issue{
				Type:    WarningIssue,
				Message: "Using ADD instead of COPY — ADD adds unneeded complexity and risk",
				Fix:     "Replace ADD with COPY for local files. Only use ADD when you need its special features (like auto-extraction of tar files)",
				Severity: "low",
				Impact:   "Potential security risks and unexpected behavior with ADD",
				References: []string{
					"https://docs.docker.com/develop/dev-best-practices/#use-copy-instead-of-add",
				},
			})
			break
		}
	}

	// Check if HEALTHCHECK is missing
	if !dockerfile.HasHealthcheck() {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "No HEALTHCHECK found",
			Fix:     "Add a HEALTHCHECK instruction, e.g.:\nHEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\\n  CMD curl -f http://localhost/ || exit 1",
			Severity: "medium",
			Impact:   "Container health status cannot be monitored",
			References: []string{
				"https://docs.docker.com/engine/reference/builder/#healthcheck",
			},
		})
	}

	// Check if USER is missing
	if !dockerfile.HasUser() {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "USER not specified — running as root",
			Fix:     "Add a non-root user and switch to it:\nRUN useradd -m myuser\nUSER myuser",
			Severity: "high",
			Impact:   "Security risk: container running with root privileges",
			References: []string{
				"https://docs.docker.com/develop/dev-best-practices/#use-non-root-users",
			},
		})
	}

	// Check for cleanup after package installations
	issues = append(issues, checkPackageCleanup(dockerfile)...)

	// Check for latest tag in base image
	if strings.HasSuffix(dockerfile.BaseImage, ":latest") || !strings.Contains(dockerfile.BaseImage, ":") {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "Using ':latest' tag or no tag specified — this is non-reproducible",
			Fix:     "Specify a fixed version tag, e.g., 'node:18.17.0' instead of 'node:latest'",
			Severity: "high",
			Impact:   "Non-reproducible builds and potential compatibility issues",
			References: []string{
				"https://docs.docker.com/develop/dev-best-practices/#use-specific-tags",
			},
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
		var fix string
		if hasAptGet {
			fix = "Add cleanup after apt-get install:\nRUN apt-get update && \\\n    apt-get install -y curl wget && \\\n    apt-get clean && \\\n    rm -rf /var/lib/apt/lists/*"
		} else if hasYum {
			fix = "Add cleanup after yum install:\nRUN yum install -y package-name && \\\n    yum clean all && \\\n    rm -rf /var/cache/yum"
		} else if hasApk {
			fix = "Add --no-cache flag to apk add:\nRUN apk add --no-cache curl wget"
		}

		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "Package installation without cleanup — adds unnecessary size",
			Fix:     fix,
			Severity: "medium",
			Impact:   "Increased image size due to package manager cache",
			References: []string{
				"https://docs.docker.com/develop/dev-best-practices/#minimize-the-number-of-layers",
			},
		})
	}

	return issues
}