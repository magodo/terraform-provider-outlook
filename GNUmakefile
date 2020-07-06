TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go')
WEBSITE_REPO=github.com/hashicorp/terraform-website
TESTTIMEOUT=180m
PKG_NAME=outlook

default: build

tools:
	@echo "==> installing required tooling..."
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/bflad/tfproviderdocs
	GO111MODULE=off go get -u github.com/katbyte/terrafmt
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$GOPATH/bin v1.24.0

build: fmtcheck
	go install

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

goimports:
	@echo "==> Fixing imports code with goimports..."
	goimports -w $(PKG_NAME)/

lint:
	./scripts/run-lint.sh

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

tflint:
	./scripts/terrafmt-acctests.sh
	go run ./linter -- ./$(PKG_NAME)/...

whitespace:
	@echo "==> Fixing source code with whitespace linter..."
	golangci-lint run ./... --no-config --disable-all --enable=whitespace --fix

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	./scripts/terrafmt-website.sh
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/
	@echo "==> Checking documentation for errors..."
	@tfproviderdocs check -provider-name=$(PKG_NAME) -require-resource-subcategory \
		-allowed-resource-subcategories-file website/allowed-subcategories

terrafmt-lint:
	@sh -c "'$(CURDIR)/scripts/terrafmt-acctests.sh'"
	@sh -c "'$(CURDIR)/scripts/terrafmt-website.sh'"

scaffold-website:
	./scripts/scaffold-website.sh

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test fmtcheck lint tools website website-lint website-test
