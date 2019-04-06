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

all: build

lint:
		$(Q)GO111MODULE=off go get -u golang.org/x/lint/golint
		$(Q)GO111MODULE=off go install golang.org/x/lint/golint
		$(Q)golint -set_exit_status $(PKGS)

chkvet:
		$(Q)test -z $$(go vet ./...)

chkfmt:
		$(Q)test -z $$(gofmt -l .)

verify: lint chkvet chkfmt test

test:
		$(Q)go test -timeout 10s ./pkg/...

fmt:
		$(Q)gofmt -w .

generate:
		$(Q)go generate ./...

build: generate verify
		$(Q)$(GOARGS) go build .
