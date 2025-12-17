.PHONY: help
help: #! Show this help message
	@echo 'Usage: make [target] ... '
	@echo ''
	@echo 'Targets:'
	@grep -h -F '#!' $(MAKEFILE_LIST) | grep -v grep | sed 's/:.*#!/:/' | column -t -s":"

.PHONY: prepare
prepare: #! Setup the project for development
	@go mod tidy

.PHONY: test
test: #! Run quality control tests
test: prepare lint unittest

.PHONY: unittest
unittest: #! Run the unit tests
	@go test ./test

.PHONY: lint
lint: #! Run helm lint against this chart
	@helm lint .

.PHONY: package
package: #! Build the helm package (e.g. ld-relay-x.y.z.tgz)
	@helm package .

.PHONY: update-golden-files
update-golden-files: #! Update unit test golden files (WARNING: Will change your local fs)
	@go test ./test -update-golden=true
