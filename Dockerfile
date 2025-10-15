# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.24 AS build-stage

ADD . /app
WORKDIR /app

COPY go.mod go.sum ./
RUN GOPROXY=https://mirrors.aliyun.com/goproxy/,direct go mod download
RUN GOPROXY=https://mirrors.aliyun.com/goproxy/,direct go get -u all

RUN CGO_ENABLED=0 GOOS=linux go build -C ./src -o /PowerWatcher

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /PowerWatcher /PowerWatcher

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/PowerWatcher"]