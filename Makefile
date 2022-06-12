bench:
	@go test -bench=. -benchmem  ./...

ut:
	@go test -race ./...

setup:
	sh ./script/setup.sh