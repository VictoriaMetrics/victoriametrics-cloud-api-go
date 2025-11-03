fmt:
	gofmt -l -w -s ./v1

vet:
	go vet ./v1/...

check-all: fmt vet golangci-lint govulncheck check-licenses

clean-checkers: remove-golangci-lint remove-govulncheck

lint: install-golangci-lint
	GOEXPERIMENT=synctest golangci-lint run

install-golangci-lint:
	which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0

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
