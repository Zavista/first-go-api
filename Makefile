build:
	@go build -o bin/gobank

run:
	@./bin/gobank

run-build:
	@go build -o bin/gobank && ./bin/gobank
	
test:
	@go test -v ./...