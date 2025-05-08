package security

import (
	"github.com/avirooppal/dock-slimscheck.git/checks"
	"github.com/avirooppal/dock-slimscheck.git/parser"
)

// RunSecurityChecks performs security checks on the Dockerfile
func RunSecurityChecks(dockerfile *parser.Dockerfile) []checks.Issue {
	var issues []checks.Issue

	// Check for root user (either explicit or implicit)
	if dockerfile.UsesRootUser() || !dockerfile.HasUser() {
		issues = append(issues, checks.Issue{
			Type:    checks.SecurityIssue,
			Message: "Container runs as root — create a non-root user",
		})
	}

	// Check for ADD with URL
	if dockerfile.HasAddWithURL() {
		issues = append(issues, checks.Issue{
			Type:    checks.SecurityIssue, 
			Message: "Using ADD with URL — risky, use curl+wget instead",
		})
	}

	// Check for unnecessary EXPOSE ports
	exposedPorts := dockerfile.GetExposedPorts()
	if len(exposedPorts) > 0 {
		// This is a simple heuristic - in a real implementation we would
		// want to check if these exposed ports are actually needed
		issues = append(issues, checks.Issue{
			Type:    checks.SecurityIssue,
			Message: "EXPOSE ports found — verify each port is necessary",
		})
	}

	// Check for COPY --chown usage
	hasCopyChown := false
	for _, inst := range dockerfile.GetInstructionsByType("COPY") {
		if inst.Arguments != "" && !hasCopyChown {
			// Check if --chown is used
			if containsFlag(inst.Arguments, "--chown") {
				hasCopyChown = true
				break
			}
		}
	}

	if !hasCopyChown && dockerfile.HasUser() {
		issues = append(issues, checks.Issue{
			Type:    checks.SecurityIssue,
			Message: "COPY without --chown flag — may cause permission issues for non-root user",
		})
	}

	// Additional security checks
	issues = append(issues, checkSecurityBestPractices(dockerfile)...)

	return issues
}

// containsFlag checks if a string contains a specific flag
func containsFlag(s, flag string) bool {
	return len(s) >= len(flag) && s[:len(flag)] == flag
}

// checkSecurityBestPractices checks for additional security best practices
func checkSecurityBestPractices(dockerfile *parser.Dockerfile) []checks.Issue {
	var issues []checks.Issue

	// Check if a specific USER is set and it's not root
	hasNonRootUser := false
	for _, inst := range dockerfile.GetInstructionsByType("USER") {
		if inst.Arguments != "root" && inst.Arguments != "0" {
			hasNonRootUser = true
			break
		}
	}

	if !hasNonRootUser {
		issues = append(issues, checks.Issue{
			Type:    checks.SecurityIssue,
			Message: "No non-root USER specified — add 'USER nonroot' or similar",
		})
	}

	// Check for HEALTHCHECK
	if !dockerfile.HasHealthcheck() {
		issues = append(issues, checks.Issue{
			Type:    checks.SecurityIssue,
			Message: "No HEALTHCHECK — add one to ensure container health monitoring",
		})
	}

	// Check for ARG usage before FROM
	hasArgBeforeFrom := false
	seenFrom := false
	
	for _, inst := range dockerfile.Instructions {
		if inst.Command == "FROM" {
			seenFrom = true
		} else if inst.Command == "ARG" && !seenFrom {
			hasArgBeforeFrom = true
			break
		}
	}
	
	if hasArgBeforeFrom {
		issues = append(issues, checks.Issue{
			Type:    checks.SecurityIssue,
			Message: "ARG used before FROM — these values persist in image history",
		})
	}

	return issues
}