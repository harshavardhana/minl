LDFLAGS := $(shell go run buildscripts/gen-ldflags.go)

all: install

checks:
	@echo "Checking deps:"
	@(env bash buildscripts/checkdeps.sh)
	@(env bash buildscripts/checkgopath.sh)

getdeps: checks
	@go get github.com/golang/lint/golint && echo "Installed golint:"
	@go get github.com/fzipp/gocyclo && echo "Installed gocyclo:"
	@go get github.com/remyoudompheng/go-misc/deadcode && echo "Installed deadcode:"
	@go get github.com/client9/misspell/cmd/misspell && echo "Installed misspell:"

# verifiers: getdeps vet fmt lint cyclo deadcode
verifiers: vet fmt lint cyclo deadcode spelling

vet:
	@echo "Running $@:"
	@GO15VENDOREXPERIMENT=1 go tool vet -all *.go
	@GO15VENDOREXPERIMENT=1 go tool vet -shadow=true *.go

spelling:
	@GO15VENDOREXPERIMENT=1 ${GOPATH}/bin/misspell *.go

fmt:
	@echo "Running $@:"
	@GO15VENDOREXPERIMENT=1 gofmt -s -l *.go

lint:
	@echo "Running $@:"
	@GO15VENDOREXPERIMENT=1 $(GOPATH)/bin/golint .

cyclo:
	@echo "Running $@:"
	@GO15VENDOREXPERIMENT=1 $(GOPATH)/bin/gocyclo -over 40 *.go

deadcode:
	@echo "Running $@:"
	@GO15VENDOREXPERIMENT=1 $(GOPATH)/bin/deadcode

build: getdeps verifiers
	@echo "Installing minl:"

test: getdeps verifiers
	@echo "Running all testing:"
	@GO15VENDOREXPERIMENT=1 go test $(GOFLAGS) ./

gomake-all: build
	@GO15VENDOREXPERIMENT=1 go build --ldflags "$(LDFLAGS)" -o $(GOPATH)/bin/minl
	@mkdir -p $(HOME)/.minl

coverage:
	@GO15VENDOREXPERIMENT=1 go test -race -coverprofile=cover.out ./
	@go tool cover -html=cover.out && echo "Visit your browser"


pkg-validate-arg-%: ;
ifndef PKG
	$(error Usage: make $(@:pkg-validate-arg-%=pkg-%) PKG=pkg_name)
endif

pkg-add: pkg-validate-arg-add
	@GO15VENDOREXPERIMENT=1 $(GOPATH)/bin/govendor add $(PKG)

pkg-update: pkg-validate-arg-update
	@GO15VENDOREXPERIMENT=1 $(GOPATH)/bin/govendor update $(PKG)

pkg-remove: pkg-validate-arg-remove
	@GO15VENDOREXPERIMENT=1 $(GOPATH)/bin/govendor remove $(PKG)

pkg-list:
	@GO15VENDOREXPERIMENT=1 $(GOPATH)/bin/govendor list

install: gomake-all

all-tests: test
	# TODO disable them for now.
	#@./tests/test-minio.sh

release: verifiers
	@MINL_RELEASE=RELEASE GO15VENDOREXPERIMENT=1 ./buildscripts/build.sh

experimental: verifiers
	@MINL_RELEASE=EXPERIMENTAL GO15VENDOREXPERIMENT=1 ./buildscripts/build.sh

clean:
	@rm -f cover.out
	@rm -f minl
	@find . -name '*.test' | xargs rm -fv
	@rm -fr release
