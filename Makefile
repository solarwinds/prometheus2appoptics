.PHONY: build clean doc test vet run docker publish

image_name := solarwinds/prometheus2appoptics
excluding_vendor := $(shell go list ./... | grep -v /vendor/)

default: build

build:
	go build -i -o prometheus2appoptics

clean:
	rm -f prometheus2appoptics

run:
	make build && ./prometheus2appoptics

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
