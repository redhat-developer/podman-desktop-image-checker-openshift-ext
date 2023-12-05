###
### Makefile Navigation
###
#
###
### Variables & Definitions
###

# Default shell `/bin/sh` has different meanings depending on the platform.
SHELL := /bin/bash
GO ?= go
GO_LDFLAGS:= $(shell if $(GO) version|grep -q gccgo ; then echo "-gccgoflags"; else echo "-ldflags"; fi)
GOCMD = CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO)
COVERAGE_PATH ?= .coverage
DESTDIR ?=
EPOCH_TEST_COMMIT ?= $(shell git merge-base $${DEST_BRANCH:-main} HEAD)
HEAD ?= HEAD
PROJECT := github.com/lstocchi/docker-openshift-analyzer
GIT_BASE_BRANCH ?= origin/main
REMOTETAGS ?= remote containers_image_openpgp
BUILDTAGS = $(REMOTETAGS)

SOURCES = $(shell find . -path './.*' -prune -o \( \( -name '*.go' -o -name '*.c' \) -a ! -name '*_test.go' \) -print)


COMMIT_NO ?= $(shell git rev-parse HEAD 2> /dev/null || true)
GIT_COMMIT ?= $(if $(shell git status --porcelain --untracked-files=no),$(call err_if_empty,COMMIT_NO)-dirty,$(COMMIT_NO))
DATE_FMT = %s
ifdef SOURCE_DATE_EPOCH
	BUILD_INFO ?= $(shell date -u -d "@$(call err_if_empty,SOURCE_DATE_EPOCH)" "+$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "+$(DATE_FMT)" 2>/dev/null || date -u "+$(DATE_FMT)")
else
	BUILD_INFO ?= $(shell date "+$(DATE_FMT)")
endif
GOFLAGS ?= -trimpath

# This must never include the 'hack' directory
export PATH := $(shell $(GO) env GOPATH)/bin:$(PATH)

CROSS_BUILD_TARGETS := \
	bin/doa.cross.linux.amd64 \
	bin/doa.cross.linux.arm64 \
	bin/doa.cross.windows.amd64 \
	bin/doa.cross.windows.arm64 \
	bin/doa.cross.darwin.amd64 \
	bin/doa.cross.darwin.arm64

CROSS_BUILD_TARGETS_ARM64 := \
	bin/doa.cross.linux.arm64 \
	bin/doa.cross.windows.arm64 \
	bin/doa.cross.darwin.arm64

CROSS_BUILD_TARGETS_AMD64 := \
	bin/doa.cross.linux.amd64 \
	bin/doa.cross.windows.amd64 \
	bin/doa.cross.darwin.amd64

# Dereference variable $(1), return value if non-empty, otherwise raise an error.
err_if_empty = $(if $(strip $($(1))),$(strip $($(1))),$(error Required variable $(1) value is undefined, whitespace, or empty))

CGO_ENABLED ?= 0
# Default to the native OS type and architecture unless otherwise specified
NATIVE_GOOS := $(shell env -u GOOS $(GO) env GOOS)
GOOS ?= $(call err_if_empty,NATIVE_GOOS)
# Default to the native architecture type
NATIVE_GOARCH := $(shell env -u GOARCH $(GO) env GOARCH)
GOARCH ?= $(NATIVE_GOARCH)
ifeq ($(call err_if_empty,GOOS),windows)
BINSFX := .exe
else
BINSFX :=
endif



###
### Primary entry-point targets
###

.PHONY: default
default: all

.PHONY: all
all: binaries

.PHONY: binaries
binaries: doa

.PHONY: vendor
vendor:
	$(GO) mod tidy -compat=1.18
	$(GO) mod vendor
	$(GO) mod verify



###
### Primary binary-build targets
###

# Make sure to warn in case we're building without the systemd buildtag.
bin/doa: $(SOURCES) go.mod go.sum
	$(GOCMD) build \
		-tags "$(BUILDTAGS)" \
		-o $@


.PHONY: doa
doa: bin/doa

###
### Secondary binary-build targets
###

# DO NOT USE: use local-cross instead
bin/doa.cross.%:
	TARGET="$*" \
	GOOS="$${TARGET%%.*}" \
	GOARCH="$${TARGET##*.}" \
	CGO_ENABLED=0 \
		$(GO) build \
		-tags '$(BUILDTAGS)' \
		-o "$@"

.PHONY: local-cross
local-cross: $(CROSS_BUILD_TARGETS) ## Cross compile binary for multiple architectures

.PHONY: local-cross-amd64
local-cross-amd64: $(CROSS_BUILD_TARGETS_AMD64) ## Cross compile binary for several AMD64 systems

.PHONY: local-cross-arm64
local-cross-arm64: $(CROSS_BUILD_TARGETS_ARM64) ## Cross compile binary for several ARM64 systems

.PHONY: cross
cross: local-cross

.PHONY: test
test:
	$(GOCMD) test \
		-v \
		-tags "$(BUILDTAGS)" \
		./...


