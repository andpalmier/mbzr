# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and make for build info
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application statically
ENV CGO_ENABLED=0
RUN make build-release

# Final stage
FROM alpine:latest

WORKDIR /

# Install CA certificates for HTTPS and create a non-root user
RUN apk add --no-cache ca-certificates && \
    adduser -D -g '' appuser

# Copy binary
COPY --from=builder /app/mbzr /mbzr

# Use non-root user
USER appuser

ENTRYPOINT ["/mbzr"]
