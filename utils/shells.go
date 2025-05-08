package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCommand runs a shell command and returns its output
func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	
	return stdout.String(), nil
}

// CommandExists checks if a command exists in the system
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// ParseDockerImageSize parses the size from docker image ls output
func ParseDockerImageSize(output string) map[string]string {
	sizes := make(map[string]string)
	
	lines := strings.Split(output, "\n")
	if len(lines) <= 1 {
		return sizes
	}
	
	// Skip header line
	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) >= 3 {
			// Image name and size
			imageName := fields[0]
			size := fields[2]
			
			// If image has a tag
			if strings.Contains(imageName, ":") {
				sizes[imageName] = size
			}
		}
	}
	
	return sizes
}

// ExtractSizeInMB converts a Docker size string to MB
func ExtractSizeInMB(sizeStr string) float64 {
	// Docker sizes are like "156MB", "1.2GB", etc.
	var value float64
	var unit string
	
	_, err := fmt.Sscanf(sizeStr, "%f%s", &value, &unit)
	if err != nil {
		return 0
	}
	
	switch strings.ToUpper(unit) {
	case "B":
		return value / 1000000
	case "KB":
		return value / 1000
	case "MB":
		return value
	case "GB":
		return value * 1000
	default:
		return 0
	}
}