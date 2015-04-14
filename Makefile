run:
	go build && ./slack

test:
	go test

lint:
	golint

format:
	gofmt .

vet:
	go vet .

imports:
	goimports .
