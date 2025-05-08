package checks

// IssueType represents the type of issue found
type IssueType string

const (
	// InfoIssue represents an informational message
	InfoIssue IssueType = "info"
	
	// WarningIssue represents a warning issue
	WarningIssue IssueType = "warning"
	
	// SecurityIssue represents a security issue
	SecurityIssue IssueType = "security"
)

// Issue represents a problem found in a Dockerfile
type Issue struct {
	Type    IssueType
	Message string
	Line    int // Optional line number
}