NAME	= datapirateschallenge
BIN 	= bin/$(NAME)

all: deps build test

deps:
	go get ./...

build: 
	go build -race -o $(BIN)

run:
	$(BIN)

fmt:
	go fmt ./...

test:
	go test 

lint:
	golint ./...

docker:
	sudo docker build -t datapirates/dockergo . 
	sudo docker run --rm -p 8080:8080 datapirates/dockergo

clean:
	rm -rf bin/ data/*.jsonl
