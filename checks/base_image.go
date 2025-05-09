package checks

import (
	"fmt"
	"strings"

	"github.com/avirooppal/dock-slimscheck/parser"
)

// Common base images and their slim alternatives
var baseImageAlternatives = map[string]string{
	"node":   "node:alpine",
	"python": "python:alpine",
	"ruby":   "ruby:alpine",
	"php":    "php:alpine",
	"nginx":  "nginx:alpine",
}

// CheckBaseImage performs checks on the base image
func CheckBaseImage(dockerfile *parser.Dockerfile) []Issue {
	var issues []Issue

	// Report the base image being used
	issues = append(issues, Issue{
		Type:    InfoIssue,
		Message: fmt.Sprintf("Base image: %s", dockerfile.BaseImage),
		Severity: "info",
		Impact:   "Base image choice affects the final image size and security posture",
		References: []string{
			"https://docs.docker.com/develop/develop-images/baseimages/",
		},
	})

	// Check for large base images
	for baseImage, alternative := range baseImageAlternatives {
		if strings.HasPrefix(dockerfile.BaseImage, baseImage+":") || dockerfile.BaseImage == baseImage {
			issues = append(issues, Issue{
				Type:    WarningIssue,
				Message: fmt.Sprintf("Using large base image (%s) — consider a smaller alternative like alpine", baseImage),
				Fix:     fmt.Sprintf("Replace '%s' with '%s' in your FROM instruction", dockerfile.BaseImage, alternative),
				Severity: "medium",
				Impact:   "Larger base images increase the final image size and potential attack surface",
				References: []string{
					"https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#use-multi-stage-builds",
					"https://hub.docker.com/_/alpine",
				},
			})
			break
		}
	}

	// Check for latest tag
	if strings.HasSuffix(dockerfile.BaseImage, ":latest") || !strings.Contains(dockerfile.BaseImage, ":") {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: "Using ':latest' tag or no tag specified — this is non-reproducible",
			Fix:     fmt.Sprintf("Specify a fixed version tag for your base image, e.g., '%s:1.2.3'", strings.Split(dockerfile.BaseImage, ":")[0]),
			Severity: "high",
			Impact:   "Non-reproducible builds can lead to unexpected behavior and security issues",
			References: []string{
				"https://docs.docker.com/develop/dev-best-practices/#use-specific-tags",
			},
		})
	}

	return issues
} 