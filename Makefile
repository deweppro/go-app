TOOLS_BIN=$(shell pwd)/.tools

install:
	go mod download
	rm -rf $(TOOLS_BIN)
	mkdir -p $(TOOLS_BIN)
	GO111MODULE=off GOBIN=$(TOOLS_BIN) go get golang.org/x/tools/cmd/cover
	GO111MODULE=off GOBIN=$(TOOLS_BIN) go get github.com/mattn/goveralls
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(TOOLS_BIN) v1.38.0
	
lint:
	$(TOOLS_BIN)/golangci-lint -v run ./...

generate:
	go generate -v ./...

build:
	go build -race -v ./...

tests:
	go test -race -v -covermode=atomic -coverprofile=coverage.out ./...
	$(TOOLS_BIN)/goveralls -coverprofile=coverage.out -repotoken $(COVERALLS_TOKEN)

pre-commite: generate lint tests

ci: install build lint tests