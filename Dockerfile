# Use the official Golang image as the base image
FROM golang:1.23

# Set environment variables
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# Set the working directory in the container
WORKDIR /app

# Copy the Go project into the container
COPY . .

# Install curl, unzip, and git for downloading additional tools
RUN apt-get update && apt-get install -y curl unzip git

# Install GolangCI-Lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.49.0

# Install SonarScanner
RUN curl -Lo /tmp/sonar-scanner.zip https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-4.7.0.2747-linux.zip
RUN unzip /tmp/sonar-scanner.zip -d /opt
RUN rm /tmp/sonar-scanner.zip
RUN ln -s /opt/sonar-scanner-4.7.0.2747-linux/bin/sonar-scanner /usr/bin/sonar-scanner
RUN chmod +x .git/hooks/pre-commit


# Add the GolangCI-Lint binary to the PATH
ENV PATH="/root/go/bin:${PATH}"

# Download Go module dependencies
RUN go mod tidy

# Run GolangCI-Lint
CMD ["golangci-lint", "run"]
