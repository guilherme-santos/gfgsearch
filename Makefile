# Pass TESTCASE when use make test or make integration test to control whick test will run
TESTCASE=
RUN_TESTCASE := $(if $(TESTCASE),-run $(TESTCASE),)

help: ## Display this help
	@ echo "Please use \`make <target>' where <target> is one of:"
	@ echo
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-10s\033[0m - %s\n", $$1, $$2}'
	@ echo

.env: ## Copy default config from .env.dist to .env (don't overwrite)
	cp .env.dist .env

run: .env build ## Run gfgsearch using .env config
	@export `cat .env | xargs`; ./gfgsearch

build:
	go build -o gfgsearch cmd/gfgsearch/main.go

test: ## Execute unit tests
	go test -race ${RUN_TESTCASE} ./...

integration-test: .env ## Execute integration tests
	@export `cat .env | xargs`; go test -race -tags integration ${RUN_TESTCASE} ./...
