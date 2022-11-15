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
test: #! Run the unit tests for this application
test: prepare lint unittest

.PHONY: unittest
unittest: #! Run the unit tests for this application
	@go test ./test

.PHONY: lint
lint: #! Run helm lint against this chart
	@helm lint

.PHONY: update-golden-files
update-golden-files: #! Update unit test golden files (WARNING: Will change your local fs)
	@go test ./test -update-golden=true
