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

s3:
	./envaws -conf ./example/configs.json -service s3 -command "./example/env.sh"

ssm:
	./envaws -conf ./example/configs.json -service ssm -command "./example/env.sh"
