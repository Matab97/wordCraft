FROM golang:1.22-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN GOOS=linux GOARCH=amd64 go build -o /main

# Expose API port
EXPOSE 8080

# Run the application
CMD ["./main"] 
