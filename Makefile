# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=restservice
CommitID=$(shell git rev-list -1 HEAD)
all: fmt test build-restservice build-crd deploy-crd create-crd deploy-restservice
fmt:
		$(GOCMD) fmt ./...
test:
		$(GOTEST) -v ./...

build-restservice:
        
		docker build -t restservice --build-arg GIT_COMMIT=${CommitID} .

build-crd:

		docker build -t crd crd/.	

deploy-crd:

		kubectl apply -f crd/kubernetes/examples/crd-deployment.yaml

create-crd:

		kubectl apply -f crd/kubernetes/examples/crd.yaml

deploy-restservice:

		kubectl apply -f crd/kubernetes/examples/example-restservice.yaml
		kubectl apply -f crd/kubernetes/examples/restservice-service.yaml						

docker-run:

		docker run -d  -p 8080:8080 restservice:latest 




