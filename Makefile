# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=restservice
all: test build
fmt:
		$(GOCMD) fmt ./...
build-Mac:

		$(GOBUILD) -o $(BINARY_NAME) -v
linux:
		GOOS=linux $(GOBUILD) -o $(BINARY_NAME) -v	
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
get-deps:
		$(GOGET)
# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker:
#need properties file to run this
		docker build -t restservice --build-arg GIT_COMMIT=$(git rev-list -1 HEAD) .
docker-run:

		docker run -d  -p 8080:8080 restservice:latest 


