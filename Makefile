run:
	go build && ./slack -token=$(token)

test:
	go test

lint:
	golint

format:
	gofmt -w=true .

vet:
	go vet .

imports:
	goimports -w=true .
