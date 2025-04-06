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

# Expose port 50052 for gRPC
EXPOSE 50052

CMD ["./main"]