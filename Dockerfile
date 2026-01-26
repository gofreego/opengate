# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o application main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/application .
COPY dev.yaml .
COPY resources/configs/routes /app/resources/configs/routes

# Set executable permission (optional as it should be inherited from build)
RUN chmod +x application

# Expose the port the application uses
EXPOSE 8083

# Define the command to run your application
CMD [ "/app/application" ]