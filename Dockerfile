# Build stage - only needed if binary doesn't exist
FROM golang:1.24-alpine as builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download || go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /tmp/server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /tmp/server ./server

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]

