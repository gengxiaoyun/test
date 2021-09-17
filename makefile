# This how we want to name the binary output
#
# use checkmake linter https://github.com/mrtazz/checkmake
# $ checkmake Makefile
#
BINARY=das
GOPATH ?= $(shell go env GOPATH)
GO111MODULE:=auto
export GO111MODULE
# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
	$(error Please set the environment variable GOPATH before running `make`)
endif
PATH := ${GOPATH}/bin:$(PATH)
GCFLAGS=-gcflags "all=-trimpath=${GOPATH}"
VERSION_TAG := $(shell git describe --tags --always)
VERSION_VERSION := $(shell git log --date=iso --pretty=format:"%cd" -1) $(VERSION_TAG)
VERSION_COMPILE := $(shell date +"%F %T %z") by $(shell go version)
VERSION_BRANCH  := $(shell git rev-parse --abbrev-ref HEAD)
VERSION_GIT_DIRTY := $(shell git diff --no-ext-diff 2>/dev/null | wc -l | awk '{print $1}')
VERSION_DEV_PATH:= $(shell pwd)
# LDFLAGS=-ldflags="-s -w -X 'github.com/XiaoMi/soar/common.Version=$(VERSION_VERSION)' -X 'github.com/XiaoMi/soar/common.Compile=$(VERSION_COMPILE)' -X 'github.com/XiaoMi/soar/common.Branch=$(VERSION_BRANCH)' -X 'github.com/XiaoMi/soar/common.GitDirty=$(VERSION_GIT_DIRTY)' -X 'github.com/XiaoMi/soar/common.DevPath=$(VERSION_DEV_PATH)'"

# These are the values we want to pass for VERSION  and BUILD
BUILD_TIME=`date +%Y%m%d%H%M`
COMMIT_VERSION=`git rev-parse HEAD`

# colors compatible setting
CRED:=$(shell tput setaf 1 2>/dev/null)
CGREEN:=$(shell tput setaf 2 2>/dev/null)
CYELLOW:=$(shell tput setaf 3 2>/dev/null)
CEND:=$(shell tput sgr0 2>/dev/null)

.PHONY: all
all: | deps clean fmt build

.PHONY: go_version_check
GO_VERSION_MIN=1.16
# Parse out the x.y or x.y.z version and output a single value x*10000+y*100+z (e.g., 1.9 is 10900)
# that allows the three components to be checked in a single comparison.
VER_TO_INT:=awk '{split(substr($$0, match ($$0, /[0-9\.]+/)), a, "."); print a[1]*10000+a[2]*100+a[3]}'
go_version_check:
	@echo "$(CGREEN)Go version check ...$(CEND)"
	@if test $(shell go version | $(VER_TO_INT) ) -lt \
	$(shell echo "$(GO_VERSION_MIN)" | $(VER_TO_INT)); \
	then printf "go version $(GO_VERSION_MIN)+ required, found: "; go version; exit 1; \
		else echo "go version check pass";	fi

# Dependency check
.PHONY: deps
deps:
	@echo "$(CGREEN)Dependency check ...$(CEND)"
	@export  GO111MODULE=on
	@export  GOPROXY=https://goproxy.cn,direct
	@go mod tidy
	@echo "download deps Success!"

# Code format
.PHONY: fmt
fmt: go_version_check
	@echo "$(CGREEN)Run gofmt on all source files ...$(CEND)"
	@echo "gofmt -l -s -w ..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		gofmt -l -s -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

# Run golang test cases
.PHONY: test
test:
	@echo "$(CGREEN)Run all test cases ...$(CEND)"
	@go test $(LDFLAGS) -timeout 10m -race ./...
	@echo "test Success!"

# Rule golang test cases with `-update` flag
.PHONY: test-update
test-update:
	@echo "$(CGREEN)Run all test cases with -update flag ...$(CEND)"
	@go test $(LDFLAGS) ./... -update
	@echo "test-update Success!"

# Code Coverage
# colorful coverage numerical >=90% GREEN, <80% RED, Other YELLOW
.PHONY: cover
cover: test
	@echo "$(CGREEN)Run test cover check ...$(CEND)"
	@go test $(LDFLAGS) -coverpkg=./... -coverprofile=coverage.data ./... | column -t
	@go tool cover -html=coverage.data -o coverage.html
	@go tool cover -func=coverage.data -o coverage.txt
	@tail -n 1 coverage.txt | awk '{sub(/%/, "", $$NF); \
		if($$NF < 80) \
			{print "$(CRED)"$$0"%$(CEND)"} \
		else if ($$NF >= 90) \
			{print "$(CGREEN)"$$0"%$(CEND)"} \
		else \
			{print "$(CYELLOW)"$$0"%$(CEND)"}}'

# Builds the project
build: fmt
	@echo "$(CGREEN)Building ...$(CEND)"
	@mkdir -p bin
	ret=0 && for d in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		b=$$(basename $${d}) ; \
		CGO_ENABLED=0 go build ${GCFLAGS} -o bin/$${b} $$d || ret=$$? ; \
	done ; exit $$ret
	@echo "build Success!"

# Installs our project: copies binaries
# install: build
# 	@echo "$(CGREEN)Install ...$(CEND)"
# 	go install ./...
# 	@echo "install Success!"

.PHONY: release
release: build
	@echo "$(CGREEN)Cross platform building for release ...$(CEND)"
	@mkdir -p release
	@for GOOS in darwin linux windows; do \
		for GOARCH in amd64; do \
			for d in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
				b=$$(basename $${d}) ; \
				echo "Building $${b}.$${GOOS}-$${GOARCH} ..."; \
				CGO_ENABLED=0 GOOS=$${GOOS} GOARCH=$${GOARCH} go build ${GCFLAGS} ${LDFLAGS} -v -o release/$${b}.$${GOOS}-$${GOARCH} $$d 2>/dev/null ; \
			done ; \
		done ;\
	done

# Cleans our projects: deletes binaries
.PHONY: clean
clean:
	@echo "$(CGREEN)Cleanup ...$(CEND)"
	go clean
	@for GOOS in darwin linux windows; do \
	    for GOARCH in 386 amd64; do \
			rm -f ${BINARY}.$${GOOS}-$${GOARCH} ;\
		done ;\
	done
	rm -f ${BINARY} 
	find . -name "*.log" -delete
	git clean -f