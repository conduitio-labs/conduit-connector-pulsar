.PHONY: build test test-integration generate install-paramgen

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/alarbada/conduit-connector-apachepulsar.version=${VERSION}'" -o conduit-connector-apachepulsar cmd/connector/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...

test-integration:
	# run required docker containers, execute integration tests, stop containers after tests
	docker compose -f test/docker-compose.yml up --quiet-pull -d
	go test $(GOTEST_FLAGS) -race ./...; ret=$$?; \
		docker compose -f test/docker-compose.yml down; \
		exit $$ret

generate:
	go generate ./...

install-paramgen:
	go install github.com/conduitio/conduit-connector-sdk/cmd/paramgen@latest
