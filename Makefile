build:
	go build -o build/release-travis ./cmd/travis
	go build -o build/release-custom ./cmd/custom

fmt:
	gofmt -w -s .
	goimports -w .

.PHONY: build
