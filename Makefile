all: test build
auto: all env
bin: test build mv-bin
	
build:
	go build .

test: 
	go tool vet .
	go test ./... -cover

mv-bin:
	mv envaws ./bin

env:
	./envaws -conf ./example/configs.json -command "./example/env.sh"
