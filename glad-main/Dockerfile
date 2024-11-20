# Copyright 2024 AboveCloud9.AI Products and Services Private Limited
# All rights reserved.
# This code may not be used, copied, modified, or distributed without explicit permission.

# Build stage
FROM golang:1.23-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Invoke 'make' command to build the Go binaries
RUN make linux-binaries

# Final stage
FROM scratch
LABEL org.opencontainers.image.authors="devops@abovecloud9.ai"

# Copy ca-certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/bin/api /api

# Copy any additional config files if needed
# COPY --from=builder /app/config /config

# Expose port (adjust as needed)
EXPOSE 8080

# TODO: uncomment once healthcheck is supported
# Add healthcheck
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#     # non-HTTP check:
#     # CMD ["/api", "health"]
#     # HTTP check: (use one of these)
#     CMD ["/api", "-http-get=http://localhost:8080/health"]

# Command to run
ENTRYPOINT ["/api"]