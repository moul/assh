regenerate:
	go install github.com/goccmack/gocc
	gocc -o ./internal/ dot.bnf 
	find . -type f -name '*.go' | xargs goimports -w

test:
	go test ./...

travis:
	make regenerate
	go build ./...
	go test ./...
	gofmt -l -s -w .
	git diff --exit-code
