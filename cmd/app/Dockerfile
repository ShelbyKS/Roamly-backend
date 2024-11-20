FROM golang:latest AS build-stage

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/app

RUN go build -o main

FROM ubuntu:latest AS build-release-stage

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /

COPY --from=build-stage /app/cmd/app/main /main

ENTRYPOINT ["/main"]