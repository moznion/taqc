.PHONY: check test lint fmt fmt-check

check: test-gen lint test fmt-check

test: test-gen
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

test-gen: clean-test-gen
	go generate ./...

clean-test-gen:
	rm -rf ./cmd/taqc/internal/tests/*_gen.go

clean:
	rm -f ./dist/taqc_*

GOLANG_CONTAINER := "golang:1.17.3-bullseye"

build:
	docker run -it --rm --env GOOS=$(GOOS) --env GOARCH=$(GOARCH) -v $(shell pwd):/taqc -w /taqc $(GOLANG_CONTAINER) \
		go build \
		-ldflags '-X "main.revision=$(shell git rev-parse HEAD)" -X "main.version=$(shell git describe --abbrev=0 --tags)"' \
		-o dist/taqc_$(GOOS)_$(GOARCH) \
		./cmd/taqc/taqc.go

