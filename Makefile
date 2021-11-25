.DEFAULT_GOAL := build

fmt:
	go fmt .

vet : fmt
	go vet

build: vet
	GOOS=darwin GOARCH=amd64 go build -o bin/github-app-auth-darwin github-app-auth.go
	GOOS=linux GOARCH=amd64 go build -o bin/github-app-auth-linux github-app-auth.go
	GOOS=windows GOARCH=amd64 go build -o bin/github-app-auth.exe github-app-auth.go
.PHONY: build

clean:
	rm -rf bin
.PHONY: clean