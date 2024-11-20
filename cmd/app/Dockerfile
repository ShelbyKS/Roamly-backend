FROM golang:latest AS build-stage

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/app

RUN go build -o main

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/cmd/app/main /main

ENTRYPOINT ["/main"]