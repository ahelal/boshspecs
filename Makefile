NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m


# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

.PHONY: all test clean deps build

all: clean deps test build

goDep:
	@echo "$(OK_COLOR)==> Installing go deps$(NO_COLOR)"
	@go get -u github.com/golang/dep/cmd/dep

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u github.com/onsi/ginkgo
# @go get -u github.com/golang/lint/golint
	@dep ensure

build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@/bin/sh -c "PKG_SRC=$(PKG_SRC) VERSION=$(VERSION) ./build/build.sh"

test:
	@echo "$(OK_COLOR)==> Testing... $(NO_COLOR)"
	@ginkgo -r -v

lint:
	@echo "$(OK_COLOR)==> Linting... $(NO_COLOR)"
	@golint `go list ./... | grep -v /vendor/`

clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	@go clean
