# Go parameters
GOCMD=go
BINARY_NAME=app
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get
GOLIST=$(GOCMD) list
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_UNIX=$(BINARY_NAME)_unix
CWD=`pwd`

all: lint test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	go test $(VERBOSE) -race `go list ./... | grep -v /vendor/`

lint:
	FILES=`go list ./... | grep -v /vendor/`
	$(GOVET) $(FILES)
	$(GOFMT) $(FILES)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker build -t promproxy .
