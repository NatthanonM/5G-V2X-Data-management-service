FROM golang:1.16.2-alpine as builder


RUN apk update && apk add --no-cache make gcc libc-dev protobuf-dev
RUN go get github.com/golang/protobuf/protoc-gen-go

ENV GO11MODULE=on

ENV SERVICE_ADDRESS=0.0.0.0:8082
ENV DATABASE_URI=
ENV DATABASE_NAME=
ENV GOOGLE_API_KEY=

WORKDIR /app
COPY Makefile /app/Makefile
ADD third_party /app/third_party
ADD cmd /app/cmd
ADD pkg /app/pkg
ADD go.mod /app/go.mod
ADD go.sum /app/go.sum
ADD api /app/api
ADD internal /app/internal

RUN go mod tidy
RUN go mod vendor

RUN make proto
RUN make bin

CMD [ "/app/bin/server-linux-amd64" ]

EXPOSE 8082
