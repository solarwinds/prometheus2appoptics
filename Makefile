.PHONY: build clean doc test vet run docker publish

image_name := solarwinds/p2l
excluding_vendor := $(shell go list ./... | grep -v /vendor/)

default: build

build:
	go build -i -o p2l

clean:
	rm -f p2l

run:
	make build && ./p2l

docker:
	GOOS=linux GOARCH=amd64 go build && docker build -t $(image_name) . && make clean && make build

publish:
	make docker && docker push $(image_name)


test:
	go test -v $(excluding_vendor)

doc:
	godoc -http=:8080 -index

vet:
	go vet ./..
