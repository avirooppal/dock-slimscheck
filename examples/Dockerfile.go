# Go application with good practices but a few issues
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Use scratch for smallest possible image
FROM scratch

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/app .

# Add CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose port - but no usage explanation
EXPOSE 8080

# No HEALTHCHECK (not easy in scratch)
# No USER specified (but scratch doesn't have users)

# Run the application
CMD ["/app/app"]