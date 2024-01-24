.PHONY: build test test-integration generate install-paramgen

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/alarbada/conduit-connector-apache-pulsar.version=${VERSION}'" -o pulsar cmd/connector/main.go

test:
	trap 'rm -f test/*.pem' EXIT
	rm -rf test/*.pem
	cd test && ./setup-tls.sh
	# run required docker containers, execute integration tests, stop containers after tests
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait 
	go test $(GOTEST_FLAGS) -v -count=1 -race .; ret=$$?; \
		docker compose -f test/docker-compose.yml down; \
		exit $$ret
	rm -f test/*.pem

generate:
	go generate ./...

install-paramgen:
	go install github.com/conduitio/conduit-connector-sdk/cmd/paramgen@latest

lint:
	golangci-lint run


up:
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait 

down:
	docker compose -f test/docker-compose.yml down
