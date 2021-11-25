.DEFAULT_GOAL := build

fmt:
	go fmt .

vet : fmt
	go vet

build: vet
	go build github-app-auth.go
.PHONY: build

clean:
	rm github-app-auth
.PHONY: clean