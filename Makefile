check: lint test

build:
	task build

tidy:
	@go mod tidy -v

test: tidy
	go test -race ./...

cover: tidy
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

cover-html: cover
	go tool cover -html=coverage.out -o coverage.html

docker:
	docker build -t docker.pkg.github.com/mach-composer/mach-composer-cli/mach:latest . --progress=plain
