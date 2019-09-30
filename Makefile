build:
	go build -o build/release-travis ./cmd/travis
	go build -o build/release-local ./cmd/local

fmt:
	gofmt -w -s .
	goimports -w .

.PHONY: build
