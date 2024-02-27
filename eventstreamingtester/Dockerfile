# Start from the latest golang base image
FROM golang:alpine

# Set the Current Working Directory inside the container
WORKDIR /eventstreamingtester

# Copy go mod and sum files
COPY . .

# Build the Go app
RUN go build -o eventstreamingtester .

# Command to run the executable
CMD ["./eventstreamingtester"]