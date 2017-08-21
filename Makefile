install:
	go install
	cp $$GOPATH/bin/escape-client $$GOPATH/bin/escape

test:
	escape test

godog-tests: install
	cd godog && godog

go-tests:
	go test -cover -v $$(go list ./... | grep -v -E 'vendor' ) | grep -v "no test files"

local-tests: go-tests godog-tests

fmt:
	find -name '*.go' | grep -v "\.escape" | grep -v vendor | xargs -n 1 go fmt
