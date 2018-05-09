include .env

.PHONY: vendor example refresh-fixtures clean-fixtures

$(GOPATH)/bin/dep:
	@go get -u github.com/golang/dep/cmd/dep

vendor: $(GOPATH)/bin/dep
	@dep ensure

example:
	@go run example/main.go

clean-fixtures:
	@-rm test/fixtures.yaml

refresh-fixtures: clean-fixtures test/fixtures.yaml

test/fixtures.yaml:
	@LINODE_TOKEN=$(LINODE_TOKEN) \
	LINODE_INSTANCE_ID=$(LINODE_INSTANCE_ID) \
	LINODE_VOLUME_ID=$(LINODE_VOLUME_ID) \
	go run test/main.go
	@sed -i "s/$(LINODE_TOKEN)/awesometokenawesometokenawesometoken/" $@

test: vendor test/fixtures.yaml
	@go test
