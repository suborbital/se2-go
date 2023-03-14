.PHONY: lint lintfix

lint:
	golangci-lint run ./...

lintfix:
	golangci-lint run ./... --fix
