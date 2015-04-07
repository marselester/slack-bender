run:
	go build && ./slack

lint:
	golint

format:
	gofmt .

vet:
	go vet .
