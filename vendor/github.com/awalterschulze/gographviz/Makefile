regenerate:
	go install github.com/goccmack/gocc
	gocc -zip -o ./internal/ dot.bnf 
	find . -type f -name '*.go' | xargs goimports -w

test:
	go test ./...

travis:
	make regenerate
	go build ./...
	go test ./...
	errcheck ./...
	gofmt -l -s -w .
	golint -set_exit_status
	git diff --exit-code
