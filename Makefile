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
test: fmt-check test-unit test-smoke test-int

.PHONY: test-unit
test-unit: fmt-check
	go test -v --tags=unit ./$(if $(PKG_NAME),linode/$(PKG_NAME),linode/...)

IP_ENV_FILE = /tmp/linode/ip_vars.env
SUBMODULE_DIR = e2e_scripts
E2E_SCRIPT = ./e2e_scripts/cloud_security_scripts/cloud_e2e_firewall/terraform-provider-linode/generate_ip_env_fw_e2e.sh

# Generate IP env file
.PHONY: generate-ip-env
generate-ip-env: $(IP_ENV_FILE)

$(IP_ENV_FILE):
	@if [ ! -d $(SUBMODULE_DIR) ]; then \
		echo "Submodule directory $(SUBMODULE_DIR) does not exist. Updating submodules..."; \
		git submodule update --init --recursive; \
	fi
	$(E2E_SCRIPT)

# Integration Test
.PHONY: test-int
test-int: fmt-check generate-ip-env
	\
	TF_VAR_ipv4_addr=$(shell grep PUBLIC_IPV4 $(IP_ENV_FILE) | cut -d '=' -f2 | tr -d '[:space:]') \
	TF_VAR_ipv6_addr=$(shell grep PUBLIC_IPV6 $(IP_ENV_FILE) | cut -d '=' -f2 | tr -d '[:space:]') \
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(if $(RUN_LONG_TESTS),$(RUN_LONG_TESTS),false) \
	set -o pipefail && go test --tags=$(if $(TEST_SUITE),$(TEST_SUITE),"integration") -v ./$(if $(PKG_NAME),linode/$(PKG_NAME),linode/...) \
	-count $(if $(COUNT),$(COUNT),1) -timeout $(if $(TIMEOUT),$(TIMEOUT),240m) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" -parallel $(if $(PARALLEL),$(PARALLEL),10) $(if $(ARGS),$(ARGS)) | awk '!/\[no test files\]/'

.PHONY: test-smoke
test-smoke: fmt-check generate-ip-env
	\
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	TF_VAR_ipv4_addr=$(shell grep PUBLIC_IPV4 $(IP_ENV_FILE) | cut -d '=' -f2 | tr -d '[:space:]') \
	TF_VAR_ipv6_addr=$(shell grep PUBLIC_IPV6 $(IP_ENV_FILE) | cut -d '=' -f2 | tr -d '[:space:]') \
	set -o pipefail && go test -v ./linode/... -run TestSmokeTests -tags=integration \
	-count $(if $(COUNT),$(COUNT),1) -timeout $(if $(TIMEOUT),$(TIMEOUT),240m) -ldflags="-X=github.com/linode/terraform-provider-linode/v2/version.ProviderVersion=acc" -parallel $(if $(PARALLEL),$(PARALLEL),10) $(if $(ARGS),$(ARGS)) | sed -e "/testing: warning: no tests to run/,+1d" -e "/\[no test files\]/d" -e "/\[no tests to run\]/d";


MARKDOWNLINT_IMG := 06kellyjac/markdownlint-cli
MARKDOWNLINT_TAG := 0.28.1

.PHONY: docs-check
docs-check:
	# markdown linter for the documents
	docker run --rm \
		-v $$(pwd):/markdown:ro \
		$(MARKDOWNLINT_IMG):$(MARKDOWNLINT_TAG) \
		--config .markdownlint.yml \
		docs

SWEEP?="tf_test,tf-test"

.PHONY: sweep
sweep:
	# sweep cleans the test infra from your account
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test -v ./$(if $(PKG_NAME),linode/$(PKG_NAME),linode/...) -sweep=$(SWEEP) $(ARGS)
