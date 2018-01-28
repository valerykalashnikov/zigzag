PACKAGES := $$(go list ./... | grep -v /vendor/)
test:
	@echo "Running tests..."
	go test $(PACKAGES)

build-server:
		go build -o bin/zigzag_server github.com/valerykalashnikov/zigzag/cmd/zigzag_server

run-server: build-server
	./bin/zigzag_server
