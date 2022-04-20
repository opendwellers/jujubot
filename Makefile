.PHONY: run

# Golang Flags
GOFLAGS ?= $(GOFLAGS:)
GO=go

mod:
	$(GO) mod download $(GOFLAGS)
	$(GO) mod tidy $(GOFLAGS)

build:
	$(GO) build -o bin/jujubot $(GOFLAGS)

run:
	$(GO) run $(GOFLAGS) $(GO_LINKER_FLAGS) *.go
