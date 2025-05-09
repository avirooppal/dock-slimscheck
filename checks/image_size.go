package checks

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/avirooppal/dock-slimscheck/parser"
)

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
					Fix:     suggestMultistagePattern(imageName),
					Severity: "high",
					Impact:   fmt.Sprintf("Large layer size (%dMB) increases image size and deployment time", sizeMB),
					References: []string{
						"https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#use-multi-stage-builds",
					},
				})
			}
			
			// If a layer grows significantly (more than 30% of previous layer)
			if prevSize > 0 && size > int64(float64(prevSize)*1.3) && size-prevSize > 50000000 {
				growthMB := (size - prevSize) / 1000000
				issues = append(issues, Issue{
					Type:    WarningIssue,
					Message: fmt.Sprintf("Layer %d grows by %dMB — check for unneeded files", i, growthMB),
					Fix:     "Review the RUN instruction for this layer and ensure all temporary files are cleaned up",
					Severity: "medium",
					Impact:   fmt.Sprintf("Layer growth of %dMB indicates potential file cleanup issues", growthMB),
					References: []string{
						"https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#minimize-the-number-of-layers",
					},
				})
			}
		}
		
		prevSize = size
	}
	
	return issues
}

// suggestMultistagePattern returns advice for using multistage builds
func suggestMultistagePattern(imageName string) string {
	// Extract base image language
	baseImageName := strings.Split(imageName, ":")[0]
	
	// Remove any registry prefix if present
	if strings.Contains(baseImageName, "/") {
		parts := strings.Split(baseImageName, "/")
		baseImageName = parts[len(parts)-1]
	}
	
	// Match and return appropriate multistage pattern
	switch {
	case strings.HasPrefix(baseImageName, "node"):
		return "Use a multistage build:\n\n# Build stage\nFROM node:slim AS build\nWORKDIR /app\nCOPY package*.json ./\nRUN npm install\nCOPY . .\nRUN npm run build\n\n# Production stage\nFROM node:alpine\nWORKDIR /app\nCOPY --from=build /app/dist ./dist\nCMD [\"node\", \"dist/main.js\"]"
	case strings.HasPrefix(baseImageName, "python"):
		return "Use a multistage build:\n\n# Build stage\nFROM python:slim AS build\nWORKDIR /app\nCOPY requirements.txt .\nRUN pip install --no-cache-dir -r requirements.txt\nCOPY . .\n\n# Production stage\nFROM python:alpine\nWORKDIR /app\nCOPY --from=build /app .\nCMD [\"python\", \"app.py\"]"
	case strings.HasPrefix(baseImageName, "openjdk"):
		return "Use a multistage build:\n\n# Build stage\nFROM openjdk:slim AS build\nWORKDIR /app\nCOPY . .\nRUN ./mvnw package\n\n# Production stage\nFROM eclipse-temurin:slim\nWORKDIR /app\nCOPY --from=build /app/target/*.jar app.jar\nCMD [\"java\", \"-jar\", \"app.jar\"]"
	case strings.HasPrefix(baseImageName, "golang"):
		return "Use a multistage build:\n\n# Build stage\nFROM golang:alpine AS build\nWORKDIR /app\nCOPY . .\nRUN go build -o main .\n\n# Production stage\nFROM alpine\nWORKDIR /app\nCOPY --from=build /app/main .\nCMD [\"./main\"]"
	default:
		return "Consider using a multistage build to reduce image size. Example structure:\n\n# Build stage\nFROM base-image AS build\n# ... build steps ...\n\n# Production stage\nFROM slim-base-image\n# ... copy artifacts and set up runtime ..."
	}
}