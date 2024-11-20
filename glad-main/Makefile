# Copyright 2024 AboveCloud9.AI Products and Services Private Limited
# All rights reserved.
# This code may not be used, copied, modified, or distributed without explicit permission.

.PHONY: all
all: build
FORCE: ;

SHELL  := env LIBRARY_ENV=$(LIBRARY_ENV) $(SHELL)
LIBRARY_ENV ?= dev

BIN_DIR = $(PWD)/bin

.PHONY: build

clean:
	rm -rf bin/*

dependencies:
	go mod download

build: dependencies build-api #build-cmd

build-api: 
	go build -tags $(LIBRARY_ENV) -o ./bin/api api/main.go

#build-cmd:
#	go build -tags $(LIBRARY_ENV) -o ./bin/search cmd/main.go

# Adding more flags for complete static linking
linux-binaries:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "$(LIBRARY_ENV) netgo" -installsuffix netgo -ldflags="-w -s" -o $(BIN_DIR)/api api/main.go
#	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "$(LIBRARY_ENV) netgo" -installsuffix netgo -ldflags="-w -s" -o $(BIN_DIR)/search cmd/main.go

ci: dependencies test	

# builds one docker image
# TODO: extend for multiple services
docker:
	docker build -t api:$(LIBRARY_ENV) .

build-mocks:
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@~/go/bin/mockgen -source=usecase/center/interface.go -destination=usecase/center/mock/center.go

test:
	go test -tags testing ./...

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done