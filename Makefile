# Pass TESTCASE when use make test or make integration test to control whick test will run
TESTCASE=
RUN_TESTCASE := $(if $(TESTCASE),-run $(TESTCASE),)

help: ## Display this help
	@ echo "Please use \`make <target>' where <target> is one of:"
	@ echo
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-10s\033[0m - %s\n", $$1, $$2}'
	@ echo

test: ## Execute unit tests
	go test -race ${RUN_TESTCASE} ./...
