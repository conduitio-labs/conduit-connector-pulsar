.PHONY: build test test-integration generate install-paramgen

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/alarbada/conduit-connector-apache-pulsar.version=${VERSION}'" -o conduit-connector-pulsar cmd/connector/main.go

test:
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait 
	go test -v -count=1 -race .; ret=$$?; \
		docker compose -f test/docker-compose.yml down && \
		exit $$ret

test-tls:
	rm -rf test/*.pem
	cd test && ./setup-tls.sh
	docker compose -f test/docker-compose-tls.yml up --quiet-pull -d --wait 
	export PULSAR_TLS=true && \
	go test -v -count=1 -run TLS -race .; ret=$$?; \
		docker compose -f test/docker-compose-tls.yml down && \
		rm -rf test/*.pem && \
		exit $$ret


test-debug:
	make test GOTEST_FLAGS="-v -count=1"

acceptance:
	go test -run Acceptance -v -count=1 .

generate:
	go generate ./...

install-paramgen:
	go install github.com/conduitio/conduit-connector-sdk/cmd/paramgen@latest

lint:
	golangci-lint run

clean:
	rm -rf test/*.pem

up:
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait 

down:
	docker compose -f test/docker-compose.yml down
