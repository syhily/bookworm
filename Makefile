.PHONY: help build test deps clean

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

gen: ## Generate templ code
	go install github.com/a-h/templ/cmd/templ@latest
	templ generate

build: deps gen ## Build executable files
	go build -gcflags=all="-l -B" --ldflags="-s -w" -buildvcs=false -o bookworm .

test: ## Run tests
	go install "github.com/rakyll/gotest@latest"
	gotest -v -coverprofile=coverage.out -covermode=atomic ./...

deps: ## Update dependencies
	go mod verify
	go mod tidy -v
	go get -u ./...

watch: ## Start a dev server
	go install github.com/romshark/templier@latest
	templier --config ./templier.yml
