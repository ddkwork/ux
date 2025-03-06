git rev-parse HEAD
go install mvdan.cc/gofumpt@latest
gofumpt -l -w .
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
pause