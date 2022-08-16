.PHONY: help
help: #! Show this help message
	@echo 'Usage: make [target] ... '
	@echo ''
	@echo 'Targets:'
	@fgrep -h '#!' $(MAKEFILE_LIST) | fgrep -v fgrep | sed -s 's/:.*#!/:/' | column -t -s":"

.PHONY: test
test: #! Run the unit tests for this application
	@go test ./test

.PHONY: update-golden-files
update-golden-files: #! Update unit test golden files (WARNING: Will change your local fs)
	@go test ./test -update-golden=true