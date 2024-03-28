COUNT?=1
PARALLEL?=10
PKG_NAME=linode/...
TIMEOUT?=240m
SWEEP?="tf_test,tf-test"
TEST_TAGS=integration

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
format:
	gofumpt -l -w .

.PHONY: fmt-check err-check imports-check vet
fmt-check:
	golangci-lint run --disable-all --enable gofumpt ./...
err-check:
	golangci-lint run --disable-all -E errcheck ./...
imports-check:
	golangci-lint run --disable-all --enable goimports ./...
vet:
	golangci-lint run --disable-all --enable govet ./...

.PHONY: test
test: fmt-check smoke-test unit-test int-test

.PHONY: unit-test
unit-test: fmt-check
	go test -v --tags=unit ./$(PKG_NAME)

.PHONY: int-test
int-test: fmt-check
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	go test --tags=$(TEST_TAGS) -v ./$(PKG_NAME) -count $(COUNT) -timeout $(TIMEOUT) -parallel=$(PARALLEL) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" $(ARGS)

.PHONY: smoke-test
smoke-test: fmt-check
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	go test -v -run smoke ./linode/... -count $(COUNT) -timeout $(TIMEOUT) -parallel=$(PARALLEL) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc"

.PHONY: opt-test
opt-test: fmt-check
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	go test --tags='integration optional' -v ./$(PKG_NAME) -count $(COUNT) -timeout $(TIMEOUT) -parallel=$(PARALLEL) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" $(ARGS)

.PHONY: long-test
long-test: fmt-check
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	go test --tags='integration long_running' -v ./$(PKG_NAME) -count $(COUNT) -timeout $(TIMEOUT) -parallel=$(PARALLEL) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" $(ARGS)

.PHONY: all-test
all-test: fmt-check
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	go test --tags='integration optional long_running' -v ./$(PKG_NAME) -count $(COUNT) -timeout $(TIMEOUT) -parallel=$(PARALLEL) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" $(ARGS)

.PHONY: docs-check
docs-check:
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
