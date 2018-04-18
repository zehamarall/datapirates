NAME	= datapirateschallenge
BIN 	= bin/$(NAME)

all: build test

# test: deps
# 	go test $(glide novendor)

# deps:
# 	go get .

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

clean:
	rm -rf bin/
