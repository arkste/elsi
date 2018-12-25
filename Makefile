all: install

build:
	go build

install:
	go install

test:
	gofmt -d -s .
	go vet ./...
	golint -set_exit_status ./...
	go test ./...