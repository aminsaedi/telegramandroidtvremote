# Stage 1: Build the application
FROM golang:1.19-alpine as builder

# Install necessary packages
RUN apk add --no-cache git

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./
COPY *.pem ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /tv-bot

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /tv-bot /tv-bot
COPY --from=builder /app/*.pem ./

# Set the binary as the entrypoint
ENTRYPOINT ["/tv-bot"]

