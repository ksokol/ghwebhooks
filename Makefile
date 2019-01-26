default: build

build:
	go build -o bin/ghwebhooks main.go

start: build
	bin/ghwebhooks
