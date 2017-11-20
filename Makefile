install:
	go install

test:
	escape run test

go-test:
	go test -cover -v $$(go list ./... | grep -v -E 'vendor' ) | grep -v "no test files"

local-test: go-test

fmt:
	find -name '*.go' | grep -v "\.escape" | grep -v vendor | xargs -n 1 go fmt

docs-build:
	escape run release -f --skip-tests --skip-deploy && cd ../escape-integration-tests && escape run release --skip-build --skip-deploy && cd -
