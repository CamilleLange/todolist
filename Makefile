project_name:=$(shell grep "module " <go.mod | sed 's;.*/;;g')

full: upd-vendor format lint deploy

$(project_name):
	go build -a -mod=vendor -o ./$(@)

build:
	$(MAKE) $(project_name)

format:
	go fmt ./...

upd-vendor:
	go mod tidy
	go mod vendor

lint: 
	golangci-lint run --allow-parallel-runners -c ./.golangci-lint.yml --fix ./...

godoc:
	command -v godoc >/dev/null 2>&1 || go get golang.org/x/tools/cmd/godoc
	echo -e "Go to http://localhost:6060/internal/$(shell head -n 1 go.mod | cut -d" " -f2)?m=all\nPress Ctrl + C to exit"
	godoc 

.PHONY: full $(project_name) deploy down build format upd-vendor test lint godoc