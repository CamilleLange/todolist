format:
	go fmt ./...

upd-vendor:
	go mod tidy
	go mod vendor

test: 
	go clean -testcache
	go test -cover ./...

lint: 
	golangci-lint run --allow-parallel-runners -c ./.golangci-lint.yml --fix ./...

.PHONY: format upd-vendor test lint 