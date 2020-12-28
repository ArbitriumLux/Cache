.PHONY: build
build: 
	mkdir artifacts
	cd artifacts
	cp ./configs/server.toml ./artifacts/
	cp -r ./static/ ./artifacts/
	go build -v ./Cache-master
	cp ./Cache-master.exe ./artifacts/
	

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build