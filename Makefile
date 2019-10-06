PKG 	    := github.com/mel3kings/lambda-template
PKG_LIST	:= $(shell go list ${PKG}/...)


# Set the version to the commit reference if building on AWS
ifdef CODEBUILD_RESOLVED_SOURCE_VERSION
	VERSION := $$(echo $(CODEBUILD_RESOLVED_SOURCE_VERSION) | cut -c 1-7)
else
	VERSION := $(shell git rev-parse HEAD | cut -c 1-7)
endif

.PHONY: setup
setup: ## Setup dependencies
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | DEP_RELEASE_TAG=v0.5.0 sh
	$(MAKE) dep

.PHONY: dep
dep: ## Check dependencies
	@dep ensure -v

.PHONY: build
build: clean ## Build locally
	GOOS=linux go build -v -o handler
	zip -r handler.zip handler config

.PHONY: deploy-lambda-test
deploy-lambda-test: build ## Deploy to lambda
	aws lambda update-function-code \
  --function-name lambda-template \
  --zip-file fileb://handler.zip

.PHONY: clean
clean: ## Remove all the temporary and build files
	rm -rf handler handler.zip

.PHONY: test
test: ## Run tests
	go test -race ${PKG_LIST}

.PHONY: display-build-version
display-build-version: ## Print the current build version
	@echo "####### Build version: $(VERSION)"

.PHONY: help
help: ## Print this
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'