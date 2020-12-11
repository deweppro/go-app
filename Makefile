lint:
	golangci-lint -v run ./...

generate:
	go generate -v ./...

tests:
	go test -v ./...

pre-commite: generate lint tests

run:
	cd ./example && go run main.go
