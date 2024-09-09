# Use the official Golang image as the base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY /src .
RUN go mod init example  \
    && go mod tidy \
    && go build -o /app/myapp ./

# Expose the port your app listens on
EXPOSE 8080

# Run the binary when starting the container
CMD ["/app/myapp"]
