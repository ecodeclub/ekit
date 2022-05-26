bench:
	@go test -bench=. -benchmem  ./...

test:
	@go test -race ./...

setup:
	sh ./script/setup.sh