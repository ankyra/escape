build:
	go build 

install: build
	go install
	cp $$GOPATH/bin/escape-client $$GOPATH/bin/escape

test:
	escape test

go-test:
	go test -cover -v $$(go list ./... | grep -v -E 'vendor' ) | grep -v "no test files"

fmt:
	find -name '*.go' | grep -v "\.escape" | grep -v vendor | xargs -n 1 go fmt
	
# Needs: sudo pip install nose-html-reporting
test-html-output: install
	nosetests -s tests --with-html
