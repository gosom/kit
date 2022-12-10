TEST_COVERAGE_THRESHOLD=80.0

default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

vet: ## runs the go vet command
	go vet ./...

test: ## runs the go test command
	go test -v -race -timeout 5m ./...

test/cover: ## runs the go test command with coverage
	go test -v -race -timeout 5m ./... -coverprofile coverage.out
	go tool cover -func coverage.out
	rm coverage.out

test/cover/html: ## runs the go test command with coverage and generates html report
	go test -v -race -timeout 5m ./... -coverprofile coverage.out
	go tool cover -func coverage.out;\
	go tool cover -html=coverage.out -o coverage.html;\
	rm coverage.out

test/cover/threshold: ## returns an error if the test coverage is below the threshold
	go test -timeout 5m ./... -coverprofile coverage.out
	coverage=$$(go tool cover -func coverage.out|grep total|grep -Eo '[0-9]+\.[0-9]+');\
	rm coverage.out;\
	passed=$$(echo "$$coverage ${TEST_COVERAGE_THRESHOLD}" | awk '{print ($$1 >= $$2)}');\
	if [ $$passed -eq 0 ]; then\
		echo "Low test coverage: $$coverage < $(TEST_COVERAGE_THRESHOLD)";\
		exit 1;\
	fi

lint: ## runs the linter
	golangci-lint run

