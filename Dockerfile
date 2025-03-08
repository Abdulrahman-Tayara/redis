# Use an official Golang runtime as a parent image
FROM golang:1.22

ENV CONFIG_FILE="default"

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

RUN go mod tidy

# Build the application
RUN go build ./cmd/server/main.go

CMD ["sh", "-c", "./main -config ${CONFIG_FILE}"]