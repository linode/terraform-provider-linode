SWEEP?="tf_test,tf-test"
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=linode

MARKDOWNLINT_IMG := 06kellyjac/markdownlint-cli
MARKDOWNLINT_TAG := 0.19.0

lint: fmtcheck
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run
	go run github.com/bflad/tfproviderlint/cmd/tfproviderlint \
		-R003=false \
		-R005=false \
		-R007=false \
		-R008=false \
		-S006=false \
		-S022=false \
		./...

changelogcheck:
	go run github.com/Charliekenney23/tf-changelog-validator/cmd/tf-changelog-validator

docscheck:
	docker run --rm \
		-v $$(pwd):/markdown:ro \
		$(MARKDOWNLINT_IMG):$(MARKDOWNLINT_TAG) \
		--config .markdownlint.yml \
		website

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

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
	go test $(TEST) -v $(TESTARGS) -timeout 120m -parallel=2 -ldflags="-X=github.com/terraform-providers/terraform-provider-linode/version.ProviderVersion=acc"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
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

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build sweep test testacc vet fmt fmtcheck errcheck test-compile website website-test
