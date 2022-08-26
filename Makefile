.PHONY:	bench
bench:
	@go test -bench=. -benchmem  ./...

.PHONY:	ut
ut:
	@go test -race ./...

.PHONY:	setup
setup:
	@sh ./script/setup.sh

.PHONY:	fmt
fmt:
	@sh ./script/goimports.sh

.PHONY:	lint
lint:
	@golangci-lint run -c .golangci.yml

.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: check
check:
	@$(MAKE) fmt
	@$(MAKE) tidy
	@$(MAKE) lint