all: test build env

build:
	go build .

test: 
	go test ./... -cover

env:
	./awsenv -command "./example/env.sh"
