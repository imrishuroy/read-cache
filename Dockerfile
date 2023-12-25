FROM golang:1.21.5-alpine3.19

# Install necessary dependencies
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy dependency files and download them
COPY go.mod .
COPY go.sum .
RUN go mod download


# Copy the entire application
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the port the application runs on
EXPOSE 8080

# Run the binary program produced by go install
ENTRYPOINT ["./main"]