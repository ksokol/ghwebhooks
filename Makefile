export GO111MODULE=on

default: build

build:
	go build -o bin/ghwebhooks *.go

start: build
	bin/ghwebhooks
