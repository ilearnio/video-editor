# install dependencies
installdeps:
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install github.com/cespare/reflex@latest

start:
	go run . serve

watch:
	reflex -d none -s -R ^pb_data/ -R ^tests/ -R _test\.go$$ -- make start

build:
	@rm -rf build
	@mkdir -p build
	cp .env* build
	go build -o build/app

test:
	go clean -testcache
	ginkgo -r --fail-fast

test-watch:
	go clean -testcache
	ginkgo watch -r --fail-fast

test-unit:
	go clean -testcache
	ginkgo -r --fail-fast

test-unit-watch:
	go clean -testcache
	ginkgo watch -r --fail-fast

.PHONY: installdeps start watch build test test-watch test-unit test-unit-watch
