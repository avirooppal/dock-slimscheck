package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Instruction represents a Dockerfile instruction
type Instruction struct {
	Command   string
	Arguments string
	Line      int
	Raw       string
}

// Dockerfile represents a parsed Dockerfile
type Dockerfile struct {
	Path         string
	Instructions []Instruction
	BaseImage    string
}

// ParseDockerfile parses a Dockerfile and returns its structure
func ParseDockerfile(path string) (*Dockerfile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open Dockerfile: %v", err)
	}
	defer file.Close()

	dockerfile := &Dockerfile{
		Path:         path,
		Instructions: []Instruction{},
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	continuationLine := ""
	
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		// Handle line continuations
		trimmedLine := strings.TrimSpace(line)
		if continuationLine != "" {
			line = continuationLine + " " + trimmedLine
			continuationLine = ""
		}
		
		// Check for line continuation
		if strings.HasSuffix(trimmedLine, "\\") {
			continuationLine = strings.TrimSuffix(trimmedLine, "\\")
			continue
		}
		
		// Skip empty lines and comments
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		
		// Parse the instruction
		parts := strings.SplitN(line, " ", 2)
		command := strings.ToUpper(parts[0])
		arguments := ""
		if len(parts) > 1 {
			arguments = parts[1]
		}
		
		instruction := Instruction{
			Command:   command,
			Arguments: arguments,
			Line:      lineNum,
			Raw:       line,
		}
		
		dockerfile.Instructions = append(dockerfile.Instructions, instruction)
		
		// Capture the base image from FROM instruction
		if command == "FROM" && dockerfile.BaseImage == "" {
			// Extract just the image name and tag
			fromArgs := strings.Fields(arguments)
			// Skip AS name if present
			for i, arg := range fromArgs {
				if strings.ToUpper(arg) == "AS" && i > 0 {
					dockerfile.BaseImage = fromArgs[0]
					break
				}
			}
			
			// If no AS found, use the first argument
			if dockerfile.BaseImage == "" && len(fromArgs) > 0 {
				dockerfile.BaseImage = fromArgs[0]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning Dockerfile: %v", err)
	}

	return dockerfile, nil
}

// GetInstructionsByType returns all instructions of a specific type
func (d *Dockerfile) GetInstructionsByType(command string) []Instruction {
	var result []Instruction
	command = strings.ToUpper(command)
	
	for _, inst := range d.Instructions {
		if inst.Command == command {
			result = append(result, inst)
		}
	}
	
	return result
}

// HasInstruction checks if the Dockerfile has a specific instruction
func (d *Dockerfile) HasInstruction(command string) bool {
	command = strings.ToUpper(command)
	
	for _, inst := range d.Instructions {
		if inst.Command == command {
			return true
		}
	}
	
	return false
}

// HasHealthcheck checks if the Dockerfile has a HEALTHCHECK instruction
func (d *Dockerfile) HasHealthcheck() bool {
	return d.HasInstruction("HEALTHCHECK")
}

// HasUser checks if the Dockerfile has a USER instruction
func (d *Dockerfile) HasUser() bool {
	return d.HasInstruction("USER")
}

// UsesRootUser checks if the Dockerfile explicitly uses root as user
func (d *Dockerfile) UsesRootUser() bool {
	for _, inst := range d.Instructions {
		if inst.Command == "USER" {
			user := strings.TrimSpace(inst.Arguments)
			return user == "root" || user == "0"
		}
	}
	
	return false
}

// HasAddWithURL checks if the Dockerfile uses ADD with a URL
func (d *Dockerfile) HasAddWithURL() bool {
	for _, inst := range d.Instructions {
		if inst.Command == "ADD" {
			// Simple pattern to match URLs
			matched, _ := regexp.MatchString(`https?://`, inst.Arguments)
			if matched {
				return true
			}
		}
	}
	
	return false
}

// HasWildcardCopy checks if the Dockerfile uses COPY with wildcards
func (d *Dockerfile) HasWildcardCopy() bool {
	for _, inst := range d.Instructions {
		if inst.Command == "COPY" {
			if strings.Contains(inst.Arguments, " . ") || strings.HasSuffix(inst.Arguments, " .") {
				return true
			}
		}
	}
	
	return false
}

// GetExposedPorts returns all ports exposed in the Dockerfile
func (d *Dockerfile) GetExposedPorts() []string {
	var ports []string
	
	for _, inst := range d.Instructions {
		if inst.Command == "EXPOSE" {
			portArgs := strings.Fields(inst.Arguments)
			ports = append(ports, portArgs...)
		}
	}
	
	return ports
}