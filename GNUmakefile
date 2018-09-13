TEST?=$$(go list ./... |grep -v 'vendor')
TESTARGS?=-race -coverprofile=profile.out -covermode=atomic
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TARGETS=darwin linux windows
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=mongodbatlas

default: build

build: fmtcheck
	go install

targets: $(TARGETS)

$(TARGETS):
	GOOS=$@ GOARCH=amd64 go build -o "dist/terraform-provider-mongodbatlas_${TRAVIS_TAG}_$@_amd64"
	zip -j dist/terraform-provider-mongodbatlas_${TRAVIS_TAG}_$@_amd64.zip dist/terraform-provider-mongodbatlas_${TRAVIS_TAG}_$@_amd64

test: fmtcheck
	TEST="$(TEST)" TESTARGS="$(TESTARGS)" sh -c "'$(CURDIR)/scripts/gotest.sh'"

testacc: fmtcheck
	TF_ACC=1 TEST="$(TEST)" TESTARGS="$(TESTARGS) -timeout 120m" sh -c "'$(CURDIR)/scripts/gotest.sh'"

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

vendor-status:
	@govendor status

vendor-fetch:
	@govendor fetch +external +missing +vendor

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
	# For testing until an official provider
	if [ ! -h $(GOPATH)/src/$(WEBSITE_REPO)/content/source/docs/providers/mongodbatlas ]; then \
		ln -s ../../../../ext/providers/mongodbatlas/website/docs/ $(GOPATH)/src/$(WEBSITE_REPO)/content/source/docs/providers/mongodbatlas; \
	fi
	if [ ! -h $(GOPATH)/src/$(WEBSITE_REPO)/content/source/layouts/mongodbatlas.erb ]; then \
		ln -s ../../../../ext/providers/mongodbatlas/website/mongodbatlas.erb $(GOPATH)/src/$(WEBSITE_REPO)/content/source/layouts/mongodbatlas.erb; \
	fi
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	# For testing until an official provider
	if [ ! -h $(GOPATH)/src/$(WEBSITE_REPO)/content/source/docs/providers/mongodbatlas ]; then \
		ln -s ../../../../ext/providers/mongodbatlas/website/docs/ $(GOPATH)/src/$(WEBSITE_REPO)/content/source/docs/providers/mongodbatlas; \
	fi
	if [ ! -h $(GOPATH)/src/$(WEBSITE_REPO)/content/source/layouts/mongodbatlas.erb ]; then \
		ln -s ../../../../ext/providers/mongodbatlas/website/mongodbatlas.erb $(GOPATH)/src/$(WEBSITE_REPO)/content/source/layouts/mongodbatlas.erb; \
	fi
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile website website-test
