GOLANGCI_LINT_VERSION := 2.12.2

fmt:
	gofmt -l -w -s ./v1

vet:
	go vet ./v1/...

check-all: fmt vet lint govulncheck check-licenses

clean-checkers: remove-golangci-lint remove-govulncheck

lint: install-golangci-lint
	golangci-lint run --build-tags 'synctest'

install-golangci-lint:
	which golangci-lint && (golangci-lint --version | grep -q $(GOLANGCI_LINT_VERSION)) || curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION)

remove-golangci-lint:
	rm -rf `which golangci-lint`

govulncheck: install-govulncheck
	govulncheck ./...

install-govulncheck:
	which govulncheck || go install golang.org/x/vuln/cmd/govulncheck@latest

remove-govulncheck:
	rm -rf `which govulncheck`

install-wwhrd:
	which wwhrd || go install github.com/frapposelli/wwhrd@latest

check-licenses: install-wwhrd
	wwhrd check -f .wwhrd.yml

test:
	go test ./v1/...
