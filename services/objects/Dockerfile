# build stage
FROM golang:1.24.2-alpine3.21 AS build-stage

WORKDIR /amdp-task/services/objects

COPY ./services/objects/go.mod ./services/objects/go.sum ./

COPY ./services/objects/internal/ internal/

COPY ./services/objects/cmd/ cmd/

COPY ./system/ ../../system/

ENV GOPROXY=https://goproxy.io,direct

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd


# final stage
FROM alpine:3.21.3

#RUN apk add --no-cache ca-certificates

WORKDIR /amdp-task/services/objects

COPY --from=build-stage /amdp-task/services/objects/app .

RUN mkdir -p /amdp-task/services/objects/logs

ENTRYPOINT [ "/amdp-task/services/objects/app" ]

