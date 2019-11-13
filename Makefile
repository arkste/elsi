all: build

build:
	@go build

install:
	@go install

commit: dep lint vet format test

lint:
	@golint -set_exit_status ./...

vet:
	@go vet ./...

format:
	@goimports -d .

test:
	@go test ./...

dep:
	@go get -v -d ./...