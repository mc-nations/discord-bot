# Use the official Golang image as the base image
FROM golang:1.21-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o nations-discord-bot

# Set the command to run the application when the container starts
CMD ["./nations-discord-bot"]
