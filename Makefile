PROJECT_NAME := terminal-calculator-tdd-go
TEST_OPTS := -covermode=atomic $(TEST_OPTS)

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: gen_all_mocks
gen_all_mocks: ## Generate all mocks
	@mockery --dir "." --output "mocks" --all