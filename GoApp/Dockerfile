# Start from the official Go image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
# GOARCH is set to amd64 to match the X86_64 architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main ./cmd/server

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy necessary directories and files
COPY --from=builder /app/static /root/static
COPY --from=builder /app/internal /root/internal

# Install CA certificates and any other necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set the timezone if needed (optional)
# ENV TZ=UTC

# Expose the port your application runs on
EXPOSE 80

# Command to run the executable
CMD ["./main"]
