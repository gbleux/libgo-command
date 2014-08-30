version := 1.0
branch:= master
create := $(sh date +%Y-%m-%d)

GOCMD = go

INSTALL ?= install
INSTALL_DATA := $(INSTALL) -m 644
INSTALL_PROGRAM := $(INSTALL) -m 755
INSTALL_DIR := $(INSTALL) -m 755 -d

MKDIR ?= $(INSTALL_DIR)
RMDIR ?= rm -vr

EXISTS_DIR ?= test -d
EXISTS_FILR ?= test -f

BASE := $(PWD)

VOLATILE_DIR := $(BASE)/build
VOLATILE_lib_DIR := $(VOLATILE_DIR)/lib
VOLATILE_bin_DIR := $(VOLATILE_DIR)/bin
DIST_DIR := $(BASE)/libgo-command-$(version)
VENDOR_DIR := $(BASE)/vendor

TESTS := test
TESTS_DIR := $(BASE)/$(TESTS)
SOURCES := src
SOURCES_DIR := $(BASE)/$(SOURCES)
EXAMPLES := src/examples
EXAMPLES_DIR := $(BASE)/$(EXAMPLES)

LDFLAGS = -X command.VERSION $(version) \
			-X command.VCS_BRANCH $(branch) \
			-X command.BUILD_DATE $(create)

AUX_files := README.md \
			LICENSE

SOURCE_dirs := command
SOURCE_targets := $(SOURCE_dirs:%=build-%)

TEST_dirs := $(SOURCE_dirs)
TEST_targets := $(TEST_dirs:%=test-%)

EXAMPLE_dirs := simple \
				playground
EXAMPLE_targets := $(EXAMPLE_dirs:%=build-example-%)

default: build

all: test

build: library examples

dist: library $(DIST_DIR)

volatile: $(VOLATILE_lib_DIR) $(VOLATILE_bin_DIR)

$(DIST_DIR):
	@$(MKDIR) "$(DIST_DIR)"

$(VOLATILE_lib_DIR):
	@$(MKDIR) "$(VOLATILE_lib_DIR)"

$(VOLATILE_bin_DIR):
	@$(MKDIR) "$(VOLATILE_bin_DIR)"

library: $(SOURCE_targets)

test: $(TEST_targets)

examples: $(EXAMPLE_targets)

$(SOURCE_targets): $(VOLATILE_lib_DIR) $(SOURCES_DIR)/$(@:build-%=%)
	@cd $(SOURCES_DIR) && \
	GOPATH="$(VOLATILE_lib_DIR):$(SOURCES_DIR):$(BASE):$(VENDOR_DIR)" \
	$(GOCMD) build \
		-ldflags "$(LDFLAGS)" \
		-o "$(VOLATILE_lib_DIR)/$(@:build-%=%)" \
		./$(@:build-%=%)

$(TEST_targets): $(TESTS_DIR)/$(@:test-%=%)
	@cd $(SOURCES_DIR) && \
	GOPATH="$(SOURCES_DIR):$(TESTS_DIR):$(BASE):$(VENDOR_DIR)" \
	$(GOCMD) test \
		-ldflags "$(LDFLAGS)" \
		./$(@:test-%=%)

$(EXAMPLE_targets): $(VOLATILE_bin_DIR) $(EXAMPLES_DIR)/$(@:build-example-%=%) library
	@cd $(EXAMPLES_DIR) && \
	GOPATH="$(VOLATILE_bin_DIR):$(VOLATILE_lib_DIR):$(EXAMPLES_DIR):$(BASE):$(VENDOR_DIR)" \
	GOBIN="$(VOLATILE_bin_DIR)" \
	$(GOCMD) build \
		-ldflags "$(LDFLAGS)" \
		-o "$(VOLATILE_bin_DIR)/$(@:build-example-%=%)" \
		./$(@:build-example-%=%)

deps:
	# install dependencies into a separate directory
	@cd $(SOURCES_DIR) && \
	GOPATH="$(VENDOR_DIR)" \
	$(GOCMD) get \
		-d \
		.

test-deps:
	@cd $(TESTS_DIR) && \
	GOPATH="$(VENDOR_DIR)" \
	$(GOCMD) get \
		-d -t \
		.

format:
	cd $(SOURCES_DIR) && \
	$(GOCMD) fmt -w .

clean:
	@$(EXISTS_DIR) "$(VOLATILE_DIR)" && $(RMDIR) $(VOLATILE_DIR) || true
	@$(EXISTS_DIR) "$(DIST_DIR)" && $(RMDIR) $(DIST_DIR) || true

.PHONY: default all
.PHONY: $(SOURCE_targets) $(EXAMPLE_targets)
.PHONY: build volatile clean dist
.PHONY: deps test-deps format
.PHONE: library examples test