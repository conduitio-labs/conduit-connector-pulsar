.PHONY: build test test-integration generate install-paramgen

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/alarbada/conduit-connector-apache-pulsar.version=${VERSION}'" -o pulsar cmd/connector/main.go

test:
	rm -rf test/*.pem
	cd test && ./setup-tls.sh
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait 
	go test $(GOTEST_FLAGS) -race .; ret=$$?; \
		rm -f test/*.pem && \
		docker compose -f test/docker-compose.yml down && \
		exit $$ret

test-debug:
	make test GOTEST_FLAGS="-v -count=1"

generate:
	go generate ./...

install-paramgen:
	go install github.com/conduitio/conduit-connector-sdk/cmd/paramgen@latest

lint:
	golangci-lint run

clean:
	rm -rf test/*.pem

up:
	cd test && ./setup-tls.sh
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait 

down:
	docker compose -f test/docker-compose.yml down
	rm -rf test/*.pem
