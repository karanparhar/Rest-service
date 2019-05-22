FROM golang:alpine AS build-evn

RUN  apk add --update \
    git

RUN mkdir -p /go/src/github.com/Rest-service/crd


COPY . /go/src/github.com/Rest-service/crd

WORKDIR /go/src/github.com/Rest-service/crd



RUN go get . && go build -o crd


FROM alpine

WORKDIR /root/crd
COPY --from=build-evn /go/src/github.com/Rest-service/crd/crd /root/crd
ENTRYPOINT [ "./crd"]