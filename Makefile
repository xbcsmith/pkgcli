PACKAGE  = github.com/xbcsmith/pkgcli
BINARY   = bin/pkgcli
COMMIT  ?= $(shell git rev-parse --short=16 HEAD)
gitversion := $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo 0.1.0-0)
VERSION ?= $(gitversion)

TOOLS    = $(CURDIR)/tools
PKGS     = $(or $(PKG),$(shell $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
GOLDFLAGS = "-X $(PACKAGE)/config.Version=$(VERSION) -X $(PACKAGE)/config.Commit=$(COMMIT)"

export GO111MODULE=on

# Allow tags to be set on command-line, but don't set them
# by default
override TAGS := $(and $(TAGS),-tags $(TAGS))

GO      = go
GOBUILD = CGO_ENABLED=0 go build -v
GOVET   = go vet
GODOC   = godoc
GOFMT   = gofmt
GOGENERATE = go generate
TIMEOUT = 15


V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: static-tests test $(BINARY) $(BINARY)-arm64 $(BINARY)-ppc64le $(BINARY)-darwin

.PHONY: release
release: update-version static-tests test $(BINARY) $(BINARY)-arm64 $(BINARY)-ppc64le $(BINARY)-darwin revert-version

.PHONY: static-tests
static-tests: fmt lint vet ## Run fmt lint and vet against all source

.PHONY: linux
linux: static-tests test $(BINARY)

.PHONY: linux-release
linux-release: update-version static-tests test $(BINARY) revert-version


.PHONY: darwin
darwin: static-tests test $(BINARY)-darwin



SOURCES = $(shell find -name vendor -prune -o -name \*.go -print)

$(BINARY): $(SOURCES); $(info $(M) building linux executable…) @ ## Build program binary
	$Q GOOS=linux GOARCH=amd64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

$(BINARY)-arm64: $(SOURCES); $(info $(M) building arm64 executable…) @ ## Build program binary for arm64
	$Q GOOS=linux GOARCH=arm64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

$(BINARY)-ppc64le: $(SOURCES); $(info $(M) building ppc64le executable…) @ ## Build program binary for ppc64le
	$Q GOOS=linux GOARCH=ppc64le $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .

$(BINARY)-darwin: $(SOURCES); $(info $(M) building darwin executable…) @ ## Build program binary
	$Q GOOS=darwin GOARCH=amd64 $(GOBUILD) $(TAGS) -ldflags $(GOLDFLAGS) -o $@ .



# Tools
GOIMPORTS = $(TOOLS)/goimports
$(GOIMPORTS): ; $(info $(M) building goimports…)
	$Q go build -o $@ golang.org/x/tools/cmd/goimports

GOLINT = $(TOOLS)/golint 
$(GOLINT): ; $(info $(M) building golint…)
	$Q go build -o $@ golang.org/x/lint/golint

GOCOVMERGE = $(TOOLS)/gocovmerge
$(GOCOVMERGE): ; $(info $(M) building gocovmerge…)
	$Q go build -o $@ github.com/wadey/gocovmerge

GOCOV = $(TOOLS)/gocov
$(GOCOV): ; $(info $(M) building gocov…)
	$Q go build -o $@ github.com/axw/gocov/gocov

GOCOVXML = $(TOOLS)/gocov-xml
$(GOCOVXML): ; $(info $(M) building gocov-xml…)
	$Q go build -o $@ github.com/AlekSi/gocov-xml

GO2XUNIT = $(TOOLS)/go2xunit
$(GO2XUNIT): ; $(info $(M) building go2xunit…)
	$Q go build -o $@ github.com/tebeka/go2xunit

GOBINDATA = $(TOOLS)/go-bindata
$(GOBINDATA): ; $(info $(M) building go-bindata…)
	@mkdir -p $(TOOLS)
	$Q go build -o $@ github.com/shuLhan/go-bindata/cmd/go-bindata

GOVERSIONINFO = $(TOOLS)/goversioninfo
$(GOVERSIONINFO): ; $(info $(M) building goversioninfo…)
	@mkdir -p $(TOOLS)
	$Q go build -o $@ github.com/josephspurrier/goversioninfo/cmd/goversioninfo


$(TOOLS)/protoc-gen-go: ; $(info $(M) building protoc-gen-go…)
	@mkdir -p $(TOOLS)
	$Q go build -o $@ github.com/golang/protobuf/protoc-gen-go

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml unit-tests functional-tests check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check tests: fmt lint ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q 2>&1 $(GO) test -timeout 20s -v $(TESTPKGS) | tee tests/tests.output
	$(GO2XUNIT) -fail -input tests/tests.output -output tests/tests.xml

COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/tests/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint test-coverage-tools ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)/coverage
	$Q for pkg in $(TESTPKGS); do \
        $(GO) test \
            -coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $$pkg | \
                    grep '^$(PACKAGE)/' | \
                    tr '\n' ',')$$pkg \
            -covermode=$(COVERAGE_MODE) \
            -coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
     done
	$Q $(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: update-version
update-version:
	$Q ./scripts/update_version.sh

.PHONY: revert-version
revert-version:
	$Q git checkout cmd/version.go

.PHONY: lint
lint: $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint change ret=1 to make lint required
	$Q ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) -set_exit_status=true $$pkg | tee /dev/stderr)" || ret=0 ; \
	 done ; exit $$ret

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

.PHONY: vet
vet: ; $(info $(M) running go vet…) @ ## Run go vet on all source files
	$(GOVET) ./...

.PHONY: test
test: ; $(info $(M) running tests…) @
	$Q go test -v ./...

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin tools vendor
	@rm -rf tests/tests.* tests/coverage.*
	@rm -rf cmd/version.bak.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)



