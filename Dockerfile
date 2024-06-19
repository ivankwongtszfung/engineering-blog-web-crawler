# Use an official Golang runtime as a parent image
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace.
ADD . /app

# Build the Go app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a Docker multi-stage build to minimize the size of the final image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binaries from the builder stage to the production image
COPY --from=builder /app/main .


# Set the entrypoint script to be executed
CMD ["./main"]