FROM golang:latest AS build-stage

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/notifier

RUN go build -o notifier

FROM ubuntu:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/cmd/notifier/notifier /notifier

ENTRYPOINT ["/notifier"]