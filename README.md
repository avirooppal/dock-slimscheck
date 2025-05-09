# Dock-SlimCheck

A Dockerfile analysis tool that helps identify potential issues, security concerns, and optimization opportunities in your Docker images.

## Features

- **Base Image Analysis**: Identifies large base images and suggests smaller alternatives
- **Best Practices Check**: Validates Dockerfile against best practices
- **Layer Size Analysis**: Analyzes layer sizes and suggests optimizations
- **Security Checks**: Identifies potential security issues (when enabled)
- **Multistage Build Suggestions**: Recommends multistage build patterns based on your base image

## Installation

### Prerequisites
- Go 1.21 or later
- Docker (optional, for layer size analysis)

<<<<<<< HEAD
### Global Installation
```bash
# Clone the repository
git clone https://github.com/avirooppal/dock-slimscheck.git
cd dock-slimscheck

# Install globally
go install

# Verify installation
dock-slimcheck --version
```

=======
>>>>>>> d1a9448fc09eb35e909d300ccd808f1d97aa9007
### Building from Source
```bash
# Clone the repository
git clone https://github.com/avirooppal/dock-slimscheck.git
cd dock-slimscheck

# Build the executable
go build -o dock-slimcheck.exe
```

## Usage

Basic usage:
```bash
<<<<<<< HEAD
dock-slimcheck ./path/to/Dockerfile
=======
dock-slimcheck.exe ./path/to/Dockerfile
>>>>>>> d1a9448fc09eb35e909d300ccd808f1d97aa9007
```

With security checks enabled:
```bash
<<<<<<< HEAD
dock-slimcheck ./path/to/Dockerfile --security
=======
dock-slimcheck.exe ./path/to/Dockerfile --security
>>>>>>> d1a9448fc09eb35e909d300ccd808f1d97aa9007
```

Show version:
```bash
<<<<<<< HEAD
dock-slimcheck --version
=======
dock-slimcheck.exe --version
>>>>>>> d1a9448fc09eb35e909d300ccd808f1d97aa9007
```

## Checks Performed

### Base Image Checks
- Identifies large base images (node, python, ruby, etc.)
- Suggests smaller alternatives (alpine, slim variants)
- Warns about using latest tags

### Best Practices
- Checks for .dockerignore when using COPY . .
- Validates ADD vs COPY usage
- Verifies HEALTHCHECK presence
- Checks for USER specification
- Validates package manager cleanup
- Warns about using latest tags

### Layer Size Analysis
- Identifies large layers (>100MB)
- Detects significant layer growth
- Suggests multistage builds when appropriate

### Security Checks (--security flag)
- Validates non-root user usage
- Checks for ADD with URL usage
- Verifies EXPOSE port necessity
- Validates COPY --chown usage
- Checks for ARG usage before FROM

## Example Output

```
[INFO] Checking Dockerfile: ./Dockerfile

[+] Base image: node:latest
[!] Using large base image (node) — consider a smaller alternative like alpine
[!] Using ':latest' tag or no tag specified — this is non-reproducible
[!] No HEALTHCHECK found
[!] USER not specified — running as root
[SECURITY] Container runs as root — create a non-root user

[✓] Check complete — 5 issues found
```

## Examples

The `examples/` directory contains various Dockerfile examples demonstrating:
- Good practices (`Dockerfile.good`)
- Security issues (`Dockerfile.security_issues`)
- Layering issues (`Dockerfile.layering_issues`)
- Language-specific examples (Python, Go, Java, etc.)
- Multistage builds (`Dockerfile.multistage`)
- Full optimized example (`Dockerfile.full_optimized`)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Version

Current version: 1.0.0
