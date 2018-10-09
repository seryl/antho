APPNAME=antho
VERBOSE?=false
VERSION=$(shell cat version.go | grep 'const Version' | awk '{print $$4}' | xargs echo)
DATESTAMP=$(date +'%Y%m%d%H%M%S')

TEST?=./...
VETARGS?= -asmdec1 -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
EXTERNAL_TOOLS=\
	github.com/golang/dep/cmd/dep \
	github.com/mitchellh/gox \
	github.com/golang/lint/golint \
	golang.org/x/tools/cmd/cover

default: test

all: test check

# bin generates the releaseable binaries
bin: generate
	@sh -c "'$(CURDIR)/scripts/build.sh'"

# bin generates the releaseable binaries for amd64 docker builds
linux: generate
	@DOCKER_MODE=1 sh -c "'$(CURDIR)/scripts/build.sh'"

docker: linux
	docker build . -t $(APPNAME):$(VERSION)

docker-ci: linux
	docker build . -t $(APPNAME):$(VERSION)-$(DATESTAMP)

# dev creates binaries for testing locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: generate
	@DEV_MODE=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# test runs the unit tests and vets the code
test: generate
	go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

# testv runs the unit tests verbosely and vets the code
testv: TESTARGS=-v
testv: test

# testacc runs acceptance tests
testacc: generate
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	go test $(TEST) -v $(TESTARGS) -timeout 45m

# testrace runs the race checker
testrace: generate
	go test -race $(TEST) $(TESTARGS)

cover:
	./scripts/coverage.sh --html

lint:
	golint $(TEST)

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go list -f '{{.Dir}}' ./... \
		| grep -v '*.bitbucket.org/pack/packapi$$' \
		| xargs go tool vet ; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Vet found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# generate runs `go generate` to build the dynamically generated
# source files.
generate:
	@go generate $(TEST)

# bootstrap the build by downloading additional tools
bootstrap:
	@for tool in $(EXTERNAL_TOOLS) ; do \
		echo "Installing $$tool" ; \
	go get $$tool; \
		done ; \
	dep ensure

clean:
	@rm -rf build/releases/* bin/*

.PHONY: bin default generate test vet bootstrap
