run:
	go build && ./slack-bender -token=$(token)

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
