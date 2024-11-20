FROM golang:1.23-alpine AS build-stage

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/notifier

RUN go build -o notifier

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/cmd/notifier/notifier /notifier

ENTRYPOINT ["/notifier"]