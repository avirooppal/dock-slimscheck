package checks

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/avirooppal/dock-slimscheck.git/parser"
)

// LargeBaseImages is a map of known large base images
var LargeBaseImages = map[string]bool{
	"node":      true,
	"python":    true,
	"ruby":      true,
	"openjdk":   true,
	"php":       true,
	"ubuntu":    true,
	"debian":    true,
	"centos":    true,
	"fedora":    true,
}

// CheckBaseImage checks base image choices
func CheckBaseImage(dockerfile *parser.Dockerfile) []Issue {
	var issues []Issue

	// First, let's show what base image is being used
	issues = append(issues, Issue{
		Type:    InfoIssue,
		Message: fmt.Sprintf("Base image: %s", dockerfile.BaseImage),
	})

	// Extract the base image name without tag
	baseImageName := strings.Split(dockerfile.BaseImage, ":")[0]
	
	// Remove any registry prefix if present
	if strings.Contains(baseImageName, "/") {
		parts := strings.Split(baseImageName, "/")
		baseImageName = parts[len(parts)-1]
	}

	// Check if using a known large base image
	isLarge := false
	for largeName := range LargeBaseImages {
		if strings.HasPrefix(baseImageName, largeName) {
			isLarge = true
			break
		}
	}

	if isLarge {
		issues = append(issues, Issue{
			Type:    WarningIssue,
			Message: fmt.Sprintf("Using large base image (%s) — consider a smaller alternative like alpine", baseImageName),
		})
	}

	return issues
}

// IsDockerAvailable checks if Docker is available in the system
func IsDockerAvailable() bool {
	cmd := exec.Command("docker", "--version")
	return cmd.Run() == nil
}

// CheckLayerSizes analyzes the layer sizes of a Docker image
func CheckLayerSizes(dockerfile *parser.Dockerfile, contextDir string) []Issue {
	var issues []Issue

	// Check if there are any FROM instructions
	fromInstructions := dockerfile.GetInstructionsByType("FROM")
	if len(fromInstructions) == 0 {
		return issues
	}

	// Try to build the image if not already built
	// This is a basic implementation and might need enhancement
	// For now, we'll just check if Docker is available and assume the image exists
	
	// Check for large layers
	largeLayerIssues := checkForLargeLayers(dockerfile.BaseImage)
	if len(largeLayerIssues) > 0 {
		issues = append(issues, largeLayerIssues...)
	}

	return issues
}

// checkForLargeLayers uses docker history to find large layers
func checkForLargeLayers(imageName string) []Issue {
	var issues []Issue
	
	// Execute docker history command
	cmd := exec.Command("docker", "history", "--human=false", "--format", "{{.Size}}", imageName)
	output, err := cmd.Output()
	if err != nil {
		// If we can't get the history, just return no issues
		return issues
	}
	
	// Parse the output to get layer sizes
	layers := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	// Convert sizes to integers and check for large layers
	var prevSize int64
	for i, layerSizeStr := range layers {
		size, err := strconv.ParseInt(layerSizeStr, 10, 64)
		if err != nil {
			continue
		}
		
		// Skip the first layer which is often the base image
		if i > 0 {
			// If a layer is significantly larger than 100MB (100000000 bytes)
			if size > 100000000 {
				sizeMB := size / 1000000
				issues = append(issues, Issue{
					Type:    WarningIssue,
					Message: fmt.Sprintf("Layer %d adds %dMB — consider using multistage builds", i, sizeMB),
				})
			}
			
			// If a layer grows significantly (more than 30% of previous layer)
			if prevSize > 0 && size > int64(float64(prevSize)*1.3) && size-prevSize > 50000000 {
				growthMB := (size - prevSize) / 1000000
				issues = append(issues, Issue{
					Type:    WarningIssue,
					Message: fmt.Sprintf("Layer %d grows by %dMB — check for unneeded files", i, growthMB),
				})
			}
		}
		
		prevSize = size
	}
	
	return issues
}

// suggestMultistagePattern returns advice for using multistage builds
func suggestMultistagePattern(dockerfile *parser.Dockerfile) string {
	// Extract base image language
	baseImageName := strings.Split(dockerfile.BaseImage, ":")[0]
	
	// Remove any registry prefix if present
	if strings.Contains(baseImageName, "/") {
		parts := strings.Split(baseImageName, "/")
		baseImageName = parts[len(parts)-1]
	}
	
	// Match and return appropriate multistage pattern
	switch {
	case strings.HasPrefix(baseImageName, "node"):
		return "Use multistage build: FROM node:slim AS build + FROM node:alpine"
	case strings.HasPrefix(baseImageName, "python"):
		return "Use multistage build: FROM python:slim AS build + FROM python:alpine"
	case strings.HasPrefix(baseImageName, "openjdk"):
		return "Use multistage build: FROM openjdk:slim AS build + FROM eclipse-temurin:slim"
	case strings.HasPrefix(baseImageName, "golang"):
		return "Use multistage build: FROM golang:alpine AS build + FROM alpine"
	default:
		return "Consider using a multistage build to reduce image size"
	}
}