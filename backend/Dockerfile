# Use the official golang image as a base image
FROM arm64v8/golang:1.22


# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 443 to the outside world
EXPOSE 8292

# Command to run the executable
CMD ["./main"]