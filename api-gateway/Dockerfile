FROM golang:1.22-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main cmd/main.go

# Expose port 8080 for HTTP
EXPOSE 8080

CMD ["./main"]