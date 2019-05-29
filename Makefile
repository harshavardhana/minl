PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
LDFLAGS := $(shell go run buildscripts/gen-ldflags.go)

TAG ?= $(USER)
BUILD_LDFLAGS := '$(LDFLAGS)'

all: build

checks:
	@echo "Checking dependencies"
	@(env bash $(PWD)/buildscripts/checkdeps.sh)

getdeps:
	@mkdir -p ${GOPATH}/bin
	@which golint 1>/dev/null || (echo "Installing golint" && go get -u golang.org/x/lint/golint)
	@which staticcheck 1>/dev/null || (echo "Installing staticcheck" && wget --quiet -O ${GOPATH}/bin/staticcheck https://github.com/dominikh/go-tools/releases/download/2019.1/staticcheck_linux_amd64 && chmod +x ${GOPATH}/bin/staticcheck)
	@which misspell 1>/dev/null || (echo "Installing misspell" && wget --quiet https://github.com/client9/misspell/releases/download/v0.3.4/misspell_0.3.4_linux_64bit.tar.gz && tar xf misspell_0.3.4_linux_64bit.tar.gz && mv misspell ${GOPATH}/bin/misspell && chmod +x ${GOPATH}/bin/misspell && rm -f misspell_0.3.4_linux_64bit.tar.gz)

verifiers: getdeps vet fmt lint staticcheck spelling

vet:
	@echo "Running $@"
	@GO111MODULE=on go vet github.com/minio/minl/...

fmt:
	@echo "Running $@"
	@GO111MODULE=on gofmt -d github.com/minio/minl/...

lint:
	@echo "Running $@"
	@GO111MODULE=on ${GOPATH}/bin/golint -set_exit_status github.com/minio/minl/...

staticcheck:
	@echo "Running $@"
	@GO111MODULE=on ${GOPATH}/bin/staticcheck github.com/minio/minl/...

spelling:
	@GO111MODULE=on ${GOPATH}/bin/misspell -locale US -error `find .`

# Builds minio, runs the verifiers then runs the tests.
check: test
test: verifiers build
	@echo "Running unit tests"
	@GO111MODULE=on go test ./... 1>/dev/null

# Builds minio locally.
build: checks
	@echo "Building minl binary to './minl'"
	@GO111MODULE=on GOFLAGS="" go build --ldflags $(BUILD_LDFLAGS) -o $(PWD)/minl 1>/dev/null

# Builds minio and installs it to $GOPATH/bin.
install: build
	@echo "Installing minl binary to '$(GOPATH)/bin/minl'"
	@mkdir -p $(GOPATH)/bin && cp -f $(PWD)/minl $(GOPATH)/bin/minl
	@echo "Installation successful. To learn more, try \"minl --help\"."

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -rvf minl
	@rm -rvf build
	@rm -rvf release
