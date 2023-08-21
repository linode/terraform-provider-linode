SWEEP?="tf_test,tf-test"
GOFMT_FILES?=$$(find . -name '*.go')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=linode/...

MARKDOWNLINT_IMG := 06kellyjac/markdownlint-cli
MARKDOWNLINT_TAG := 0.19.0

ACCTEST_COUNT?=1
ACCTEST_PARALLELISM?=20
ACCTEST_TIMEOUT?=240m
RUN_LONG_TESTS?="false"

tooldeps:
	go generate -tags tools tools/tools.go

lint: fmtcheck
	golangci-lint run --timeout 15m0s
	tfproviderlint \
		-AT001=false \
		-S006=false \
		-R018=false \
		-R019=false \
		./...

docscheck:
	docker run --rm \
		-v $$(pwd):/markdown:ro \
		$(MARKDOWNLINT_IMG):$(MARKDOWNLINT_TAG) \
		--config .markdownlint.yml \
		website

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test -v ./$(PKG_NAME) -sweep=$(SWEEP) $(SWEEPARGS)

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -parallel=2 -timeout=30s

testacc: fmtcheck
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	go test --tags=integration -v ./$(PKG_NAME) $(TESTARGS) -count $(ACCTEST_COUNT) -timeout $(ACCTEST_TIMEOUT) -parallel=$(ACCTEST_PARALLELISM) -ldflags="-X=github.com/linode/terraform-provider-linode/version.ProviderVersion=acc"

smoketest: fmtcheck
	TF_ACC=1 \
	LINODE_API_VERSION="v4beta" \
	RUN_LONG_TESTS=$(RUN_LONG_TESTS) \
	go test -v -run smoke ./linode/... -count $(ACCTEST_COUNT) -timeout $(ACCTEST_TIMEOUT) -parallel=$(ACCTEST_PARALLELISM) -ldflags="-X=github.com/linode/terraform-provider-linode/version.ProviderVersion=acc"

unittest:
	go test -v --tags=unit ./linode/...
vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)
	gofumpt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS) -timeout 120m -parallel=2


imports:
	goimports -w $(GOFMT_FILES)

.PHONY: build sweep test testacc vet fmt fmtcheck errcheck test-compile
