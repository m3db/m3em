SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(SELF_DIR)/.ci/common.mk

SHELL=/bin/bash -o pipefail

html_report := coverage.html
test := .ci/test-cover.sh
test_ci_integration := .ci/test-integration.sh
convert-test-data := .ci/convert-test-data.sh
coverfile := cover.out
coverage_xml := coverage.xml
junit_xml := junit.xml
test_log := test.log
lint_check := .ci/lint.sh
m3em_package := github.com/m3db/m3em
mockgen_package := github.com/golang/mock/mockgen
vendor_prefix := vendor
gopath_prefix := $(GOPATH)/src
license_dir := .ci/uber-licence
license_node_modules := $(license_dir)/node_modules
auto_gen := .ci/auto-gen.sh
mocks_output_dir := generated/mocks/mocks
mocks_rules_dir := generated/mocks

BUILD := $(abspath ./out)
LINUX_AMD64_ENV := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
VENDOR_ENV := GO15VENDOREXPERIMENT=1

SERVICES :=  \
	m3em_agent \

setup:
	mkdir -p $(BUILD)

define SERVICE_RULES
$(SERVICE): setup
	@echo Building $(SERVICE)
	$(VENDOR_ENV) go build -o $(BUILD)/$(SERVICE) ./services/$(SERVICE)/.

prod-$(SERVICE):
	$(LINUX_AMD64_ENV) make $(SERVICE)
endef

services: $(SERVICES)
prod-services:
	$(LINUX_AMD64_ENV) make services

$(foreach SERVICE,$(SERVICES),$(eval $(SERVICE_RULES)))

lint:
	@which golint > /dev/null || go get -u github.com/golang/lint/golint
	$(VENDOR_ENV) $(lint_check)

test-internal:
	@which go-junit-report > /dev/null || go get -u github.com/sectioneight/go-junit-report
	@$(VENDOR_ENV) $(test) $(coverfile) | tee $(test_log)

test-xml: test-internal
	go-junit-report < $(test_log) > $(junit_xml)
	gocov convert $(coverfile) | gocov-xml > $(coverage_xml)
	@$(convert-test-data) $(coverage_xml)
	@rm $(coverfile) &> /dev/null

test: test-internal
	gocov convert $(coverfile) | gocov report

testhtml: test-internal
	gocov convert $(coverfile) | gocov-html > $(html_report) && open $(html_report)
	@rm -f $(test_log) &> /dev/null

test-ci-unit: test-internal
	@which goveralls > /dev/null || go get -u -f github.com/mattn/goveralls
	goveralls -coverprofile=$(coverfile) -service=travis-ci || echo -e "\x1b[31mCoveralls failed\x1b[m"

test-ci-integration:
	@$(VENDOR_ENV) $(test_ci_integration)

mock-gen: install-mockgen install-license-bin
	@echo Generating mocks
	PACKAGE=$(m3em_package) $(auto_gen) $(mocks_output_dir) $(mocks_rules_dir)

install-mockgen: install-vendor
	@echo Installing mockgen
	rm -rf $(gopath_prefix)/$(mockgen_package) && \
	cp -r $(vendor_prefix)/$(mockgen_package) $(gopath_prefix)/$(mockgen_package) && \
	go install $(mockgen_package)

install-license-bin: install-vendor
	@echo Installing node modules
	git submodule update --init --recursive
	[ -d $(license_node_modules) ] || (cd $(license_dir) && npm install)

clean:
	echo Cleaning build artifacts...
	go clean
	rm -rf $(BUILD)
	@rm -f *.html *.xml *.out *.test
	echo

.DEFAULT_GOAL := test
.PHONY: test test-xml test-internal testhtml clean
