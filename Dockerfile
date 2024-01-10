# Part 1: Builder
# Start from the latest golang base image
FROM golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set environment variables for cross-compilation
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go app
RUN make build

# Part 2: Final Server
# Start from a minimal base image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/server .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/server", "runserver"]
