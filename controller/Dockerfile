FROM golang:alpine AS build-evn

RUN mkdir -p /go/src/github.com/Rest-service/crd


COPY . /go/src/github.com/Rest-service/controller

WORKDIR /go/src/github.com/Rest-service/controller



RUN go build -o crd


FROM alpine

WORKDIR /root/crd
COPY --from=build-evn /go/src/github.com/Rest-service/controller/controller /root/controller
ENTRYPOINT [ "./crd"]