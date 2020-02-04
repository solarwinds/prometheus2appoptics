.PHONY: build clean doc test vet run docker publish

image_name := solarwinds/prometheus2appoptics
app_name := prometheus2appoptics

default: build

build:
	go build -i -ldflags="-s -w" -o $(app_name)

clean:
	rm -f $(app_name)

run:
	make build && ./$(app_name)

docker:
	docker build -t solarwinds/$(app_name):latest .

publish:
	make docker && docker push $(image_name)

test:
	go test -v ./...

doc:
	godoc -http=:8080 -index

vet:
	go vet ./...

release:
	make build
	git tag -a $(shell ./$(app_name) -version)