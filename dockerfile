# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install dependencies
COPY . .
RUN go mod download

# Build source files
RUN go build -o server

# Deployment
FROM alpine:latest
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/server .

# Expose the application port
ENV PORT=8888
EXPOSE 8888

# Run the server
CMD ["./server"]