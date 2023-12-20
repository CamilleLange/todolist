# env
project_name:=$(shell grep "module " <go.mod | sed 's;.*/;;g')

format:
	go fmt ./...

upd-vendor:
	go mod tidy
	go mod vendor

test: upd-vendor 
	go clean -testcache
	GIN_MODE=release go test -timeout 1m -cover $$(go list ./... | grep -v test)

lint: 
	golangci-lint run --allow-parallel-runners -c ./.golangci-lint.yml --fix ./...

godoc:
	command -v godoc >/dev/null 2>&1 || go get golang.org/x/tools/cmd/godoc
	echo -e "Go to http://localhost:6060/pkg/$(shell head -n 1 go.mod | cut -d" " -f2)?m=all\nPress Ctrl + C to exit"
	godoc 

.PHONY: format upd-vendor test lint godoc
