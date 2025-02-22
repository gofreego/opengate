# Use a minimal base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from your local machine to the container
COPY application ./
COPY dev.yml ./

RUN chmod +x ./application
# Expose any necessary ports (if your Go binary is a server)
EXPOSE 8000
# Command to run your Go binary
ENTRYPOINT ["./application"]