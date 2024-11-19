FROM golang:latest AS build-stage

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

# RUN go mod tidy
RUN go mod download

COPY . .

WORKDIR /app/cmd/consumer

RUN go build -o main

FROM ubuntu:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/cmd/consumer/main /main

ENTRYPOINT ["/main"]