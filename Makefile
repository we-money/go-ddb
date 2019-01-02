.PHONY: default get codetest fmt lint vet secure

default: fmt codetest

get:
	go get -t -v ./...
	go get -u golang.org/x/lint/golint
	go get -u github.com/securego/gosec/cmd/gosec/...

codetest: lint vet secure

fmt:
	go fmt ./...

lint:
	@echo golint ./...
	@OUTPUT=`golint ./... 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "golint errors:"; \
		echo "$$OUTPUT"; \
		exit 1; \
	fi

vet:
	go tool vet -all .

secure:
	gosec -quiet ./...
