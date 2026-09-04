# Dev targets for MGRS-Go. Tools install into $(GOBIN) when missing.
GO ?= go
GOBIN ?= $(shell $(GO) env GOPATH)/bin
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOFUMPT ?= $(GOBIN)/gofumpt
GOIMPORTS ?= $(GOBIN)/goimports
STATICCHECK ?= $(GOBIN)/staticcheck

GOLANGCI_LINT_VERSION ?= v2.13.2
STATICCHECK_VERSION ?= v0.8.1
GOFUMPT_VERSION ?= v0.11.0
GOIMPORTS_VERSION ?= v0.49.0

.PHONY: help tools fmt vet lint staticcheck test test-race oracle bench fuzz ci clean

help:
	@printf '%s\n' \
		'Targets:' \
		'  tools         Install golangci-lint, gofumpt, goimports, staticcheck' \
		'  fmt           Format with gofumpt and goimports' \
		'  vet           Run go vet' \
		'  lint          Run golangci-lint' \
		'  staticcheck   Run staticcheck' \
		'  test          Run go test ./...' \
		'  test-race     Run tests with the race detector' \
		'  oracle        Oracle subset (PROJ / GeoConvert / goldens)' \
		'  bench         Bench smoke for EncodeTo' \
		'  fuzz          Short bounded fuzz' \
		'  ci            Local stand-in for CI checks' \
		'  clean         Remove local build artifacts'

tools: $(GOLANGCI_LINT) $(GOFUMPT) $(GOIMPORTS) $(STATICCHECK)

$(GOLANGCI_LINT):
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(GOFUMPT):
	$(GO) install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

$(GOIMPORTS):
	$(GO) install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

$(STATICCHECK):
	$(GO) install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)

fmt: $(GOFUMPT) $(GOIMPORTS)
	$(GOFUMPT) -l -w .
	$(GOIMPORTS) -local github.com/Quad4-Software/MGRS-Go -w .

vet:
	$(GO) vet ./...

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

test:
	$(GO) test ./... -count=1

test-race:
	$(GO) test ./... -race -count=1

oracle:
	$(GO) test ./mgrs -count=1 -run 'Proj|GeoLib|Golden'

bench:
	$(GO) test ./mgrs -bench=BenchmarkEncodeTo_preallocated -benchmem -count=1 -benchtime=100ms

fuzz:
	$(GO) test ./mgrs -fuzz=FuzzEncodeDecodeIEEE -fuzztime=5s

ci: vet lint staticcheck test

clean:
	rm -f coverage.out cover.out cpu.out mem.out mgrs.test example/example
