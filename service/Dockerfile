FROM golang:alpine AS build-evn

RUN mkdir -p /go/src/github.com/Rest-service


COPY . /go/src/github.com/Rest-service/service

WORKDIR /go/src/github.com/Rest-service/service



RUN go build -o restservice


FROM alpine

WORKDIR /root/restservice
COPY --from=build-evn /go/src/github.com/Rest-service/service/restservice /root/restservice
ARG GIT_COMMIT

ENV GitCommit=$GIT_COMMIT
LABEL git_commit=$GIT_COMMIT
EXPOSE 8080
ENTRYPOINT [ "./restservice"]


