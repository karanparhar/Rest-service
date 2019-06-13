# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=restservice
CommitID=$(shell git rev-list -1 HEAD)
all: fmt test build-restservice build-crd deploy-crd
fmt:
		$(GOCMD) fmt ./service/...
test:
		$(GOTEST) -v ./service/...

build-restservice:
        
		docker build -t restservice --build-arg GIT_COMMIT=${CommitID} service/.

build-crd:

		docker build -t controller controller/.

deploy-crd:

		kubectl apply -f controller/kubernetes/examples/crd-deployment.yaml

create-crd:

		kubectl apply -f controller/kubernetes/examples/crd.yaml

deploy-restservice:

		kubectl apply -f controller/kubernetes/examples/example-restservice.yaml
		kubectl apply -f controller/kubernetes/examples/restservice-service.yaml						

docker-run:

		docker run -d  -p 8080:8080 restservice:latest 




