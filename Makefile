build:
	go build -o build/release ./cmd/release

fmt:
	gofmt -w -s .
	goimports -w .

.PHONY: build
