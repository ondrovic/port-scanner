FROM golang:1.22.5-alpine AS builder
# Set the working directory inside the container
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

# Command to run the executable
ENTRYPOINT ["./port-scanner"]
# ENTRYPOINT ["sh"]