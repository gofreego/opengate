# Build UI stage
FROM node:20-alpine AS ui-builder

WORKDIR /app/ui

# Copy package files
COPY ui/package*.json ./
RUN npm install

# Copy source and build
COPY ui/ .
RUN npm run build

# Build Go stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Copy the built UI
COPY --from=ui-builder /app/ui/dist /app/ui/dist
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
EXPOSE 8080

# Define the command to run your application
CMD [ "/app/application" ]