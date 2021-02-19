# Rest Service

This is sample rest api service with single endpoint written in go, to be used in [sample-controller](git clone https://github.com/karanparhar/Rest-service.git) project.

## Quick start

### Prerequisites
- [go](https://golang.org/dl/) version v1.10+
- minikube 

### Steps to run

```
$ mkdir $GOPATH/src/github.com/
$ cd $GOPATH/src/github.com/
$ git clone https://github.com/karanparhar/Rest-service.git
$ cd Rest-service
$ make all
$ make create-crd
$ make deploy-restservice
```

## check status
```
$ curl -X GET $(minikube service restservice --url)/api/status
```

### Note it will build rest service and crd and deploy in local kubernetes cluster 

### if role binding issue came kubectl create clusterrolebinding serviceaccounts-cluster-admin --clusterrole=cluster-admin --group=system:serviceaccounts
 
