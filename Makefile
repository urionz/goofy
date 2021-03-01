all: test lint

test:
	go test ./...

test-cover:
	go test -coverpkg=./...

test-cover-html:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golint ./...

.PHONY: test test-cover test-cover-html

