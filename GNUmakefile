SWEEP?="tf_test,tf-test"
GOFMT_FILES?=$$(find . -name '*.go')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=linode

MARKDOWNLINT_IMG := 06kellyjac/markdownlint-cli
MARKDOWNLINT_TAG := 0.19.0

ACCTEST_COUNT?=1
ACCTEST_PARALLELISM?=20
ACCTEST_POLL_MS?=1000
ACCTEST_TIMEOUT?=240m

tooldeps:
	go generate -tags tools tools/tools.go

lint: fmtcheck
	golangci-lint run
	tfproviderlint \
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
	go test ./$(PKG_NAME) -v -sweep=$(SWEEP) $(SWEEPARGS)

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
	LINODE_EVENT_POLL_MS=$(ACCTEST_POLL_MS) \
	go test ./$(PKG_NAME) -v $(TESTARGS) -count $(ACCTEST_COUNT) -timeout $(ACCTEST_TIMEOUT) -parallel=$(ACCTEST_PARALLELISM) -ldflags="-X=github.com/linode/terraform-provider-linode/version.ProviderVersion=acc"

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

.PHONY: build sweep test testacc vet fmt fmtcheck errcheck test-compile
