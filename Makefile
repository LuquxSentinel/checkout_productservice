build:
	@go build -o ./bin/productservice

run: build
	@./bin/productservice

test:
	@./...