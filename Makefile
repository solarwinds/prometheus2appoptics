.PHONY: build clean doc test vet run

excluding_vendor := $(shell go list ./... | grep -v /vendor/)

default: build

build:
	go build -i -o p2l

clean:
	rm -f p2l

run:
	make build && ./p2l

test:
	go test -v $(excluding_vendor)

doc:
	godoc -http=:8080 -index

vet:
	go vet ./..
