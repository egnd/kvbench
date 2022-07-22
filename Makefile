#!make

MAKEFLAGS += --always-make

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

%:
	@:

########################################################################################################################

size:
	sudo chown -R $$(whoami) ./
	du -h --max-depth=1 .data

gen:
	rm -rf vendor
	go generate ./...


benchmarks: gen
	go test -benchmem -bench .