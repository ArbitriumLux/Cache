.PHONY: build
build: 
	mkdir artifacts
	cd artifacts
	cp -r ./configs/ ./artifacts/
	cp -r ./static/ ./artifacts/
	go build -v main.go
	cp main ./artifacts/
	

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build