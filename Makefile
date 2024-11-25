COMMIT     :=$(shell git rev-parse HEAD)
YAML_FILES :=$(shell find . ! -path "./vendor/*" ! -path "./api/*" -type f -regex ".*y*ml" -print)

all: help

.PHONY: version
version: ## Prints the current version
	@echo $(COMMIT)

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

.PHONY: upgrade
upgrade: clean ## Upgrades all dependencies 
	go get -u ./...
	go mod tidy
	go mod vendor

.PHONY: test
test: tidy ## Runs unit tests
	go test -count=1 -race -covermode=atomic -coverprofile=cover.out ./...

.PHONY: lint
lint: lint-go lint-yaml ## Lints the entire project 
	@echo "Completed Go and YAML lints"

.PHONY: lint-go
lint-go: ## Lints the entire project using go 
	golangci-lint -c .golangci.yaml run

.PHONY: lint-yaml
lint-yaml: ## Runs yamllint on all yaml files (brew install yamllint)
	yamllint -c .yamllint $(YAML_FILES)

.PHONY: vulncheck
vulncheck: ## Checks for source vulnerabilities
	govulncheck -test ./...

.PHONY: tag
tag: ## Creates release tag 
	git tag -s -m "version bump to $(VERSION)" $(VERSION)
	git push origin $(VERSION)

.PHONY: tagless
tagless: ## Delete the current release tag 
	git tag -d $(VERSION)
	git push --delete origin $(VERSION)

.PHONY: clean
clean: ## Cleans bin and temp directories
	go clean
	rm -fr ./vendor
	rm -fr ./bin

.PHONY: api
api: ## Generates the API documentation
	oapi-codegen -config api/airthings/config.yaml api/airthings/openapi.yaml

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
