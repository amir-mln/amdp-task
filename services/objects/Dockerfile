# build stage
FROM golang:1.24.2-alpine3.21 AS build-stage

WORKDIR /amdp-task

COPY go.mod go.sum ./

COPY internal/ ./internal/

COPY cmd/ ./cmd/

ENV GOPROXY=https://goproxy.io,direct

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd


# final stage
FROM alpine:3.21.3

RUN apk add --no-cache ca-certificates

WORKDIR /amdp-task

COPY --from=build-stage /amdp-task/app .

RUN mkdir -p /amdp-task/logs

EXPOSE 3000

ENTRYPOINT [ "/amdp-task/app" ]

