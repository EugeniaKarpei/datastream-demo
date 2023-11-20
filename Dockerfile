# Start from the official Go image to create a build artifact
FROM golang:1.21.4 as builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy the go.mod and go.sum file and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Copy the data directory from the previous stage
COPY --from=builder /app/data ./data

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Start a new stage from scratch for the running container
FROM alpine:latest  

# Install ca-certificates in case you need to call HTTPS endpoints
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .
# Copy the data directory from the previous stage
COPY --from=builder /app/data ./data

# Command to run the executable
CMD ["./main"]
