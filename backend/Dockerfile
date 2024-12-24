# Build stage using Alpine
FROM golang:1.23.4-alpine as builder

# Install necessary build tools
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency installation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY .. .

# Build the Go application
RUN go build -o backend .

# Runtime stage using Alpine
FROM alpine:latest

# Install necessary runtime tools
RUN apk add ca-certificates curl

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/backend .

# Copy the env file
COPY --from=builder /app/.env .

# Ensure the binary is executable
RUN chmod +x ./backend

# Expose the port your service listens on
EXPOSE 8000

# Command to run the application
ENTRYPOINT ["./backend"]