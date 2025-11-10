watch:
	@air

build:
	echo "Building the binary in the bin dir"
	@go build -o ./bin/main ./cmd/api
