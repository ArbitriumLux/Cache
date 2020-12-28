.PHONY: build
build: 
	mkdir artifacts
	cd artifacts
	cp ./configs/server.toml ./artifacts/
	cp -r ./static/ ./artifacts/
	go build -v main.go
	cp main.exe ./artifacts/
	

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build