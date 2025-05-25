default: test

.PHONY: test
test:
	go test -v -race -coverprofile=cover.out -covermode=atomic ./...
	go tool cover -html=cover.out

.PHONY: bench
bench:
	go test ./... -run=NOTEST -bench . -benchmem

.PHONY: gen
gen:
	cd ./internal/stream && ./gen.sh
	cd ./collection/skipmap && ./gen.sh
	cd ./collection/skipset && ./gen.sh

.PHONY: license
license:
	license-eye -c .licenserc.yaml header fix