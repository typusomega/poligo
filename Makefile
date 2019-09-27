# kernel-style V=1 build verbosity
ifeq ("$(origin V)", "command line")
       BUILD_VERBOSE = $(V)
endif

ifeq ($(BUILD_VERBOSE),1)
       Q =
else
       Q = @
endif

PKGS = $(shell go list ./...)

export CGO_ENABLED:=0

all: verify build

lint: fmt
		$(Q)echo "linting...."
		$(Q)GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
		$(Q)golangci-lint run -E gofmt -E golint -E goconst -E gocritic -E golint -E gosec -E maligned -E nakedret -E prealloc -E unconvert -E gocyclo -E scopelint -E goimports
		$(Q)echo linting OK

test:
		$(Q)echo "unit testing...."
		$(Q)go test ./pkg/...

verify: test lint

clean:
		$(Q)rm -rf build

fmt:
		$(Q)echo "fixing imports and format...."
		$(Q)goimports -w .

build: verify
		$(Q)$(GOARGS) go build ./pkg/policy
