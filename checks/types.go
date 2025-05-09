package checks

// IssueType represents the type of issue found
type IssueType int

const (
	WarningIssue IssueType = iota
	SecurityIssue
	InfoIssue
)

// Issue represents a problem found in the Dockerfile
type Issue struct {
	Type        IssueType
	Message     string
	Fix         string    // Detailed fix suggestion
	Severity    string    // "low", "medium", "high"
	Impact      string    // Description of the impact
	References  []string  // Links to relevant documentation
}