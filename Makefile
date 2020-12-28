.PHONY: build
build: 
	mkdir artifacts
	cd artifacts
	cp ./configs/config.toml ./artifacts/
	cp -r ./static/ ./artifacts/
	go build -v ./Avito
	cp ./Avito.exe ./artifacts/
	

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build