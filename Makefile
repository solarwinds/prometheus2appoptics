.PHONY: build clean doc test vet run docker publish

image_name := solarwinds/prometheus2appoptics
app_name := prometheus2appoptics
excluding_vendor := $(shell go list ./... | grep -v /vendor/)

default: build

build:
	go build -i -o $(app_name)

clean:
	rm -f $(app_name)

run:
	make build && ./$(app_name)

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

release:
	make build
	git tag -a $(shell ./$(app_name) -version)