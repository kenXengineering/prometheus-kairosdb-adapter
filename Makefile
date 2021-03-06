DOCKERHUB_USERNAME ?= chosenken
IMAGE_NAME ?= prometheus-kairosdb-adapter

GOTAGS ?= prometheus-kairosdb-adapter
GOFILES ?= $(shell go list ./... | grep -v /vendor/)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

VERSION ?= latest

all: build

dep:
	dep ensure -v
	dep prune -v
	go get -u -v github.com/goreleaser/goreleaser

buildit:
	go build -o prometheus-kairosdb-adapter main.go

build: staticcheck gosimple vet buildit

install: staticcheck gosimple
	go install

run:
	go run main.go

test:
	go test $(shell go list ./... | grep -v /vendor/)

image:
	docker build -t $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(VERSION) .

push: image
	docker push $(DOCKERHUB_USERNAME)/$(IMAGE_NAME):$(VERSION)

format:
	@echo ">> Running go fmt"
	@go fmt $(GOFILES)

vet:
	@echo ">> Running go vet"
	@go vet $(GOFILES); if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

staticcheck:
	@echo ">> running staticcheck"
	@staticcheck $(GOFILES)

gosimple:
	@echo ">> running gosimple"
	@gosimple $(GOFILES)

tools:
	@echo ">> installing some extra tools"
	@go get -u -v honnef.co/go/tools/...

publish-test:
	goreleaser --skip-publish --skip-validate --rm-dist

publish:
	goreleaser --skip-validate --rm-dist

