install:
	go install

test:
	escape test

go-tests:
	go test -cover -v $$(go list ./... | grep -v -E 'vendor' ) | grep -v "no test files"

local-tests: go-tests 

fmt:
	find -name '*.go' | grep -v "\.escape" | grep -v vendor | xargs -n 1 go fmt
