# Build stage
FROM golang:1.20 AS builder

WORKDIR /ikhbodi

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /ikhbodi/main .

# Expose the port your application runs on
EXPOSE 8080

# Run the binary
CMD ["./main"]