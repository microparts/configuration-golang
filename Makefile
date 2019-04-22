IMAGE = golang-pkg-test
TEST_CONTAINER = docker run --rm -i --name golang-config $(IMAGE)

deps: ## Get the dependencies
	@go mod vendor

race: ## Run data race detector
	@go test -race ./...

tests: ## Run application unit tests with coverage and generate global code coverage report
	go test -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

covercli: ## Generate code coverage report
	@go tool cover -func=coverage.out

coverhtml: ## Generate global code coverage report in HTML
	@go tool cover -html=coverage.out

coverage: tests coverhtml

coverage_cli: tests covercli

image_test:
	@docker build -f Dockerfile -t $(IMAGE) .

tests_in_docker: image_test ## Testing code with unit tests in docker container
	$(TEST_CONTAINER) make coverage_cli

race_in_docker:
	@$(TEST_CONTAINER) make race