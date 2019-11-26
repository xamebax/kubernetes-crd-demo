.PHONY: e2e

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

generate-code:
	${ROOT_DIR}/hack/update-codegen.sh

verify:
	${ROOT_DIR}/hack/verify-codegen.sh
