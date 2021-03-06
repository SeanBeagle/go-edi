# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

# ...not implemented
# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

COPY *.go ./
ADD data ./data

RUN go env -w GO111MODULE=auto
RUN go build -o /main

EXPOSE 8080

CMD [ "/main" ]
