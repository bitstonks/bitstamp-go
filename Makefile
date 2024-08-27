gofmt-check:
	gofmt -d .

gofmt:
	gofmt -w .

goimports-install:
	go install golang.org/x/tools/cmd/goimports@latest

goimports-check: goimports-install
	goimports -d .

goimports: goimports-install
	goimports -w .

govulncheck-install:
	go install golang.org/x/vuln/cmd/govulncheck@latest

security: govulncheck-install
	govulncheck ./...

upgrade:
	go get -u ./...

format: gofmt goimports

lint: gofmt-check goimports-check
