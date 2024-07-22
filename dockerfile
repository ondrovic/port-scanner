# Build stage
FROM golang:1.22.5-alpine AS builder
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Remove the .env file
RUN rm -f .env

# Build the application
RUN go build -o port-scanner ./cmd/scanner

# Final stage
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/port-scanner /port-scanner

# Command to run the executable
ENTRYPOINT ["/port-scanner"]