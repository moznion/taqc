.PHONY: check test lint fmt fmt-check

check: gen lint test fmt-check

test: gen
	go test ./... -race -v -coverprofile="coverage.txt" -covermode=atomic

lint:
	go vet ./...
	staticcheck ./...

fmt:
	gofmt -w -s *.go
	goimports -w *.go

fmt-check:
	goimports -l *.go | grep [^*][.]go$$; \
		EXIT_CODE=$$?; \
		if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi

gen:
	go generate ./...
