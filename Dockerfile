# Dockerfile

# This Dockerfile is used to build a full scall golang project with multiple files and folders with vendor cache

# Set the base image
FROM golang:1.20.5-bullseye as builder

# Maintainer
LABEL maintainer="Sohel Ahmed Jony <sohelahmedjony@gmail.com>"
ENV GO111MODULE=on

# Set the working directory
WORKDIR /go/src/app

# Copy the source from the current directory to the working directory
COPY . .

# Install dependencies
RUN go get -d -v ./...
RUN go mod vendor
RUN go mod vendor

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 5669
EXPOSE 4001

# Set environment variables
ENV API_PORT 5669
ENV P2P_PORT 4001
# ENV GIN_MODE release

# Run the executable
CMD ["./main"]




#FROM golang:1.20.5-alpine3.18 as builder


# RUN apk --update add git


# WORKDIR /app


# RUN go mod tidy

