SWEEP?="tf_test,tf-test"
PKG_NAME=linode/...

MARKDOWNLINT_IMG := 06kellyjac/markdownlint-cli
MARKDOWNLINT_TAG := 0.28.1

ACCTEST_COUNT?=1
ACCTEST_PARALLELISM?=20
ACCTEST_TIMEOUT?=240m
RUN_LONG_TESTS?="false"

.PHONY: build sweep test testacc vet fmt fmtcheck errcheck test-compile

tooldeps:
	go generate -tags tools tools/tools.go

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

docscheck:
	docker run --rm \
		-v $$(pwd):/markdown:ro \
		$(MARKDOWNLINT_IMG):$(MARKDOWNLINT_TAG) \
		--config .markdownlint.yml \
		docs

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test -v ./$(PKG_NAME) -sweep=$(SWEEP) $(SWEEPARGS)

default: build

build: fmtcheck
	# trying to copy .goreleaser.yaml
	go build -a -ldflags '-s -extldflags "-static"'

clean:
	rm -f terraform-provider-linode

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -parallel=2 -timeout=30s

testacc: fmtcheck
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	go test --tags=integration -v ./$(PKG_NAME) -count $(ACCTEST_COUNT) -timeout $(ACCTEST_TIMEOUT) -parallel=$(ACCTEST_PARALLELISM) -ldflags="-X=github.com/linode/terraform-provider-linode/version.ProviderVersion=acc" $(TESTARGS)

smoketest: fmtcheck
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	go test -v -run smoke ./linode/... -count $(ACCTEST_COUNT) -timeout $(ACCTEST_TIMEOUT) -parallel=$(ACCTEST_PARALLELISM) -ldflags="-X=github.com/linode/terraform-provider-linode/version.ProviderVersion=acc"

unittest:
	go test -v --tags=unit ./linode/...

vet:
	golangci-lint run --disable-all --enable govet ./...

fmt:
	gofumpt -w -l .

imports:
	golangci-lint run --disable-all --enable goimports ./...

fmtcheck:
	golangci-lint run --disable-all --enable gofumpt ./...

errcheck:
	golangci-lint run --disable-all -E errcheck ./...

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS) -timeout 120m -parallel=2
