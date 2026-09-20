GO ?= go

.PHONY: help test integration vet examples check

help:
	@printf '%s\n' \
		'make test     Run the Go test suite.' \
		'make integration Run native Foundation Models integration tests.' \
		'make vet      Run go vet.' \
		'make examples Build integration examples into examples/.' \
		'make check    Run tests, vet, and build examples.'

test:
	$(GO) test ./...

integration:
	FMGO_INTEGRATION=1 $(GO) test -count=1 ./...

vet:
	$(GO) vet ./...

examples:
	mkdir -p examples
	$(GO) build -o examples/respond ./cmd/respond
	$(GO) build -o examples/server ./cmd/server
	$(GO) build -o examples/stream ./cmd/stream
	$(GO) build -o examples/structured ./cmd/structured

check: test vet examples
