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
	docker build -t docker.pkg.github.com/labd/mach-composer/mach:latest . --progress=plain


update-deps:
	go get -u github.com/mach-composer/mach-composer-plugin-aws@main
	go get -u github.com/mach-composer/mach-composer-plugin-azure@main
	go get -u github.com/mach-composer/mach-composer-plugin-amplience@main
	go get -u github.com/mach-composer/mach-composer-plugin-commercetools@main
	go get -u github.com/mach-composer/mach-composer-plugin-contentful@main
	go get -u github.com/mach-composer/mach-composer-plugin-sentry@main
	go get -u github.com/mach-composer/mach-composer-plugin-helpers@main
	go get -u github.com/mach-composer/mach-composer-plugin-sdk@main
