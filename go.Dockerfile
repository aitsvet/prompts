# Use an official Go runtime as the parent image
FROM golang:latest AS build-env

# Set the working directory in the container to /app
WORKDIR /app

# Copy local code into the container at /app
COPY . .

# Get go modules dependencies
RUN go get -d -v ./...

# Compile the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/shard/main.go

# Use a separate image to keep intermediate container clean
FROM alpine:latest

# Copy the executable from the build environment onto the smaller final image
COPY --from=build-env /app/main .

# The application will be run using cmd/shard/main as its entry point
ENTRYPOINT ["./main"]
