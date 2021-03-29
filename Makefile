lint:
	golangci-lint -v run ./...

generate:
	go generate -v ./...

tests:
	go test -race -v ./...

pre-commite: generate lint tests

