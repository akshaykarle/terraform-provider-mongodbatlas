WEBSITE_REPO=github.com/hashicorp/terraform-website
TEST?=$$(go list ./... |grep -v 'vendor')
TESTARGS?=-race -coverprofile=profile.out -covermode=atomic
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TARGETS=darwin linux windows
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=mongodbatlas

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v -parallel 20 $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

websitefmtcheck:
	@sh -c "'$(CURDIR)/scripts/websitefmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@GOGC=30 golangci-lint run ./$(PKG_NAME)

tools:
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

targets: $(TARGETS)

$(TARGETS):
	GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build -o "dist/terraform-provider-mongodbatlas_${TRAVIS_TAG}_$@_amd64"
	zip -j dist/terraform-provider-mongodbatlas_${TRAVIS_TAG}_$@_amd64.zip dist/terraform-provider-mongodbatlas_${TRAVIS_TAG}_$@_amd64

.PHONY: build test testacc fmt fmtcheck lint tools test-compile website website-lint website-test targets darwin linux windows
