# Testing `dock-slimcheck`

This document provides information about testing the `dock-slimcheck` tool with various Dockerfile examples.

## Test Files

The `examples/` directory contains several Dockerfiles designed to test different aspects of the tool:

1. **Dockerfile.sample** - Basic example with common issues
2. **Dockerfile.good** - Well-optimized example following best practices
3. **Dockerfile.python** - Python application with several best practice issues
4. **Dockerfile.java** - Java application with multistage build but other issues
5. **Dockerfile.go** - Go application with good practices but some issues
6. **Dockerfile.nginx** - Nginx example with ARG before FROM and other issues
7. **Dockerfile.full_optimized** - Fully optimized Dockerfile with all best practices
8. **Dockerfile.security_issues** - Dockerfile with multiple security issues
9. **Dockerfile.multistage** - Multistage build example for a Go application
10. **Dockerfile.layering_issues** - Example with layer optimization issues

## Running Tests

### Basic Check

To run a basic check on any example:

```bash
./dock-slimcheck examples/Dockerfile.sample
```

### Security Check

To include security checks:

```bash
./dock-slimcheck examples/Dockerfile.sample --security
```

### Testing All Examples

Test all example Dockerfiles at once:

```bash
# Linux/macOS
for file in examples/Dockerfile.*; do
  echo "Checking $file"
  ./dock-slimcheck "$file"
  echo "----------------"
done

# Windows PowerShell
Get-ChildItem "examples/Dockerfile.*" | ForEach-Object {
  Write-Host "Checking $_"
  ./dock-slimcheck $_
  Write-Host "----------------"
}
```

## Expected Results

- **Dockerfile.sample**: Multiple issues including large base image, no .dockerignore, no USER/HEALTHCHECK
- **Dockerfile.good**: Should have minimal or no issues
- **Dockerfile.python**: Large base image, no cleanup, missing USER/HEALTHCHECK
- **Dockerfile.java**: Large base image, still running as root despite multistage
- **Dockerfile.go**: Missing HEALTHCHECK, EXPOSE without explanation
- **Dockerfile.nginx**: ARG before FROM, ADD with URL
- **Dockerfile.full_optimized**: Should have minimal or no issues
- **Dockerfile.security_issues**: Multiple security issues (port exposure, root user, ADD with URL)
- **Dockerfile.multistage**: Good practices but missing HEALTHCHECK
- **Dockerfile.layering_issues**: Multiple small RUN statements, inefficient layer ordering

## Advanced Testing

To test layer size analysis, you'll need:

1. Docker installed and running
2. Images built from these Dockerfiles

```bash
# Build an image from an example
docker build -t test-image -f examples/Dockerfile.sample .

# Then analyze with dock-slimcheck
./dock-slimcheck examples/Dockerfile.sample
```