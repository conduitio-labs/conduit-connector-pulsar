VERSION=$(shell git describe --tags --dirty --always)

.PHONY: build
build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-pulsar.version=${VERSION}'" -o conduit-connector-pulsar cmd/connector/main.go

.PHONY: test
test:
	docker compose -f test/docker-compose.yml up pulsar --quiet-pull -d --wait 
	go test -v -count=1 -race .; ret=$$?; \
		docker compose -f test/docker-compose.yml down && \
		exit $$ret

.PHONY: test-tls
test-tls:
	docker compose -f test/docker-compose.yml up pulsar-tls --quiet-pull -d --wait 
	export PULSAR_TLS=true && \
	go test -v -count=1 -run TLS -race .; ret=$$?; \
		docker compose -f test/docker-compose.yml down && \
		exit $$ret

.PHONY: test-debug
test-debug:
	make test GOTEST_FLAGS="-v -count=1"

.PHONY: acceptance
acceptance:
	go test -run Acceptance -v -count=1 .

.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: install-tools
install-tools:
	@echo Installing tools from tools.go
	@go list -e -f '{{ join .Imports "\n" }}' tools.go | xargs -I % go list -f "%@{{.Module.Version}}" % | xargs -tI % go install %
	@go mod tidy

.PHONY: up
up:
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait 

.PHONY: down
down:
	docker compose -f test/docker-compose.yml down --remove-orphans
