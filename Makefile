.PHONY: vet format deps-tidy check-git-clean test build ci-test docker-build

vet:
	go vet .

test: vet
	go test ./...

format:
	gofmt -w .

deps-tidy:
	go mod tidy

check-git-clean:
	git diff
	git diff --quiet

build: ./output/multimedia-sync

./output/multimedia-sync:
	go build -o $@ .

docker-build:
	docker build .

ci-test: vet format deps-tidy check-git-clean test
