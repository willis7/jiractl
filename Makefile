VERSION?="0.1"
TEST?=./...
GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')

default: test

# bin generates the releaseable binaries for Jiractl
bin: fmtcheck generate
	@JCTL_RELEASE=1 @sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing Jiractl locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: fmtcheck generate
	go install -mod=vendor .

# test runs the unit tests
test: fmtcheck generate
	go list -mod=vendor $(TEST) | xargs -t -n4 go test $(TESTARGS) -mod=vendor -timeout=2m -parallel=4

# testrace runs the race checker
testrace: fmtcheck generate
	go test -mod=vendor -race $(TEST) $(TESTARGS)

cover:
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

# generate runs `go generate` to build the dynamically generated
# source files
generate: vendor
	GOFLAGS=-mod=vendor go generate ./...

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vendor:
	go mod vendor

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

.PHONY: bin cover default dev fmt fmtcheck generate quickdev test-compile test testacc testrace vendor-status website website-test vendor
