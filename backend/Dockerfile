# Start from a small, secure base image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download the Go module dependencies
RUN go mod download

# Build argument to specify which component to build
ARG BUILD_TARGET=main

# Build the Go binary based on BUILD_TARGET
RUN if [ "$BUILD_TARGET" = "partner-processor" ]; then \
        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/partner_applications_processor/main.go; \
    else \
        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go; \
    fi
# build for remote debug
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -gcflags="all=-N -l" -o main ./cmd/main.go

# Create a minimal production image
FROM golang:1.22-alpine
# Production image for remote debug
#FROM golang:1.21-alpine

# It's essential to regularly update the packages within the image to include security patches
RUN apk update && apk upgrade

# Reduce image size
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

# Avoid running code as a root user
RUN adduser -D appuser
USER appuser

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app/main .

# Expose the port that the application listens on
EXPOSE 8080

# remote debug app
#RUN go install github.com/go-delve/delve/cmd/dlv@latest
#EXPOSE 2345

# Run the binary when the container starts
CMD ["./main"]

# Run in debug mode
#CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "--log", "exec", "./main"]
