COUNT?=1
PARALLEL?=10
PKG_NAME=linode/...
TIMEOUT?=240m
RUN_LONG_TESTS?="false"
SWEEP?="tf_test,tf-test"

MARKDOWNLINT_IMG := 06kellyjac/markdownlint-cli
MARKDOWNLINT_TAG := 0.28.1

.PHONY: default
default: build

.PHONY: build
build: format
	# trying to copy .goreleaser.yaml
	go build -a -ldflags '-s -extldflags "-static"'

.PHONY: clean
clean:
	rm -f terraform-provider-linode

.PHONY: lint
lint:
	# remove two disabled linters when their errors are addressed
	golangci-lint run \
		--disable gosimple \
		--disable staticcheck \
		--timeout 15m0s
	tfproviderlint \
		-AT001=false \
		-AT004=false \
		-S006=false \
		-R018=false \
		-R019=false \
		./...

.PHONY: deps
deps:
	go generate -tags tools tools/tools.go

.PHONY: format
format: fmt vet errcheck imports

.PHONY: fmt vet errcheck imports
fmt:
	golangci-lint run --disable-all --enable gofumpt ./...
vet:
	golangci-lint run --disable-all --enable govet ./...
errcheck:
	golangci-lint run --disable-all -E errcheck ./...
imports:
	golangci-lint run --disable-all --enable goimports ./...

.PHONY: test
test: format smoke-test unit-test int-test

.PHONY: unit-test
unit-test:
	go test -v --tags=unit ./$(PKG_NAME)

.PHONY: int-test
int-test: format
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	go test --tags=integration -v ./$(PKG_NAME) -count $(COUNT) -timeout $(TIMEOUT) -parallel=$(PARALLEL) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" $(ARGS)

.PHONY: smoke-test
smoke-test: format
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	go test -v -run smoke ./linode/... -count $(COUNT) -timeout $(TIMEOUT) -parallel=$(PARALLEL) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc"

.PHONY: docscheck
docscheck:
	# markdown linter for the documents
	docker run --rm \
		-v $$(pwd):/markdown:ro \
		$(MARKDOWNLINT_IMG):$(MARKDOWNLINT_TAG) \
		--config .markdownlint.yml \
		docs

.PHONY: sweep
sweep:
	# sweep cleans the test infra from your account
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test -v ./$(PKG_NAME) -sweep=$(SWEEP) $(ARGS)
