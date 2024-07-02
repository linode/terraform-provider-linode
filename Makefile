COUNT?=1
PARALLEL?=10
PKG_NAME=linode/...
TIMEOUT?=240m
RUN_LONG_TESTS?=False
SWEEP?="tf_test,tf-test"
TEST_TAGS="integration"

MARKDOWNLINT_IMG := 06kellyjac/markdownlint-cli
MARKDOWNLINT_TAG := 0.28.1

IP_ENV_FILE := /tmp/linode/ip_vars.env

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
	go test -v --tags=unit ./$(PKG_NAME) | grep -v "\[no test files\]"

.PHONY: int-test
int-test: fmt-check generate-ip-env-fw-e2e include-env
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	TF_VAR_ipv4_addr=${PUBLIC_IPV4} \
	TF_VAR_ipv6_addr=${PUBLIC_IPV6} \
	go test --tags="$(TEST_TAGS)" -v ./$(PKG_NAME) -count $(COUNT) -timeout $(TIMEOUT) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" -parallel=$(PARALLEL) $(ARGS) | grep -v "\[no test files\]"

.PHONY: include-env
include-env: $(IP_ENV_FILE)
-include $(IP_ENV_FILE)

generate-ip-env-fw-e2e: $(IP_ENV_FILE)

$(IP_ENV_FILE):
	# Generate env file for E2E cloud firewall
	. ./e2e_scripts/cloud_security_scripts/cloud_e2e_firewall/terraform-provider-linode/generate_ip_env_fw_e2e.sh || touch $(IP_ENV_FILE)

.PHONY: smoke-test
smoke-test: fmt-check
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	go test -v -run smoke ./linode/... -count $(COUNT) -timeout $(TIMEOUT) -parallel=$(PARALLEL) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" | grep -v "\[no test files\]"

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
