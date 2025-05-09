# Dock-SlimCheck

A Dockerfile analysis tool that helps identify potential issues, security concerns, and optimization opportunities in your Docker images. Think of it as having a senior DevOps engineer looking over your shoulder, providing detailed suggestions for improvement.

## Features

* **Base Image Analysis**: Identifies large base images and suggests smaller alternatives
* **Best Practices Check**: Validates Dockerfile against best practices
* **Layer Size Analysis**: Analyzes layer sizes and suggests optimizations
* **Security Checks**: Identifies potential security issues (when enabled)
* **Multistage Build Suggestions**: Recommends multistage build patterns based on your base image
* **Detailed Fix Suggestions**: Not just problems, but specific solutions with code examples

## Installation

### Prerequisites

* Go 1.21 or later
* Docker (optional, for layer size analysis)

### Global Installation

```bash
# Clone the repository
git clone https://github.com/avirooppal/dock-slimscheck.git
cd dock-slimscheck

# Install globally
go install

# Verify installation
dock-slimscheck --version
```

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
dock-slimscheck ./path/to/Dockerfile
```

With security checks enabled:

```bash
dock-slimscheck ./path/to/Dockerfile --security
```

Show version:

```bash
dock-slimscheck --version
```

## Checks Performed

### Base Image Checks

* Identifies large base images (node, python, ruby, etc.)
* Suggests smaller alternatives (alpine, slim variants)
* Warns about using `latest` tags
* Provides specific version recommendations

### Best Practices

* Checks for `.dockerignore` when using `COPY . .`
* Validates `ADD` vs `COPY` usage
* Verifies `HEALTHCHECK` presence
* Checks for `USER` specification
* Validates package manager cleanup
* Warns about using `latest` tags

### Layer Size Analysis

* Identifies large layers (>100MB)
* Detects significant layer growth
* Suggests multistage builds when appropriate
* Provides language-specific multistage build examples

### Security Checks (`--security` flag)

* Validates non-root user usage
* Checks for `ADD` with URL usage
* Verifies `EXPOSE` port necessity
* Validates `COPY --chown` usage
* Checks for `ARG` usage before `FROM`

## Example Output

```bash
[INFO] Checking Dockerfile: ./Dockerfile

[+] Base image: node:latest
  → Severity: info
  → Impact: Base image choice affects the final image size and security posture
  → References:
    - https://docs.docker.com/develop/develop-images/baseimages/

[!] Using large base image (node) — consider a smaller alternative like alpine
  → Severity: medium
  → Impact: Larger base images increase the final image size and potential attack surface
  → Fix:
    Replace 'node:latest' with 'node:alpine' in your FROM instruction
  → References:
    - https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#use-multi-stage-builds
    - https://hub.docker.com/_/alpine

[!] No HEALTHCHECK found
  → Severity: medium
  → Impact: Container health status cannot be monitored
  → Fix:
    Add a HEALTHCHECK instruction, e.g.:
    HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
      CMD curl -f http://localhost/ || exit 1
  → References:
    - https://docs.docker.com/engine/reference/builder/#healthcheck

[✓] Check complete — 3 issues found
```

## Examples

The `examples/` directory contains various Dockerfile examples demonstrating:

* Good practices (`Dockerfile.good`)
* Security issues (`Dockerfile.security_issues`)
* Layering issues (`Dockerfile.layering_issues`)
* Language-specific examples (Python, Go, Java, etc.)
* Multistage builds (`Dockerfile.multistage`)
* Full optimized example (`Dockerfile.full_optimized`)

## Common Issues and Fixes

### Large Base Images
```dockerfile
# Before
FROM node:latest

# After
FROM node:18.17.0-alpine
```

### Missing Health Check
```dockerfile
# Before
CMD ["npm", "start"]

# After
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
CMD ["npm", "start"]
```

### Package Installation Cleanup
```dockerfile
# Before
RUN apt-get update && apt-get install -y curl

# After
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

### Multistage Build Example
```dockerfile
# Build stage
FROM node:slim AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Production stage
FROM node:alpine
WORKDIR /app
COPY --from=build /app/dist ./dist
CMD ["node", "dist/main.js"]
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. Here are some areas where we'd love to see contributions:

* Additional checks for common issues
* More language-specific optimizations
* Enhanced security checks
* CI/CD integration examples
* Documentation improvements

## Version

Current version: **1.0.0**

## License

MIT License - feel free to use this tool in your projects!

## Acknowledgments

* Docker best practices documentation
* Alpine Linux project
* The open-source community

