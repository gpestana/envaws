all: test build
auto: all env

build:
	go build .

test: 
	go tool vet .
	go test ./... -cover

env:
	./envaws -conf ./example/configs.json -command "./example/env.sh"
