run: build
	./bin/logmanager

build:
	go build -o ./bin/logmanager main.go

test:
	go test -v ./...

setup:
	mkdir -p datastores/

.PHONY: run build test
