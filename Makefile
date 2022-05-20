build:
	go build -o mach-composer

lint:
	staticcheck ./...

format:
	go fmt ./...

test:
	go test -v ./...

coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... ./...
	go tool cover -func=coverage.txt

docker_aws:
	docker build -t docker.pkg.github.com/labd/mach-composer/mach:latest-aws --build-arg PROVIDER=aws . --progress=plain

docker_azure:
	docker build -t docker.pkg.github.com/labd/mach-composer/mach:latest-azure --build-arg PROVIDER=azure . --progress=plain

docker:
	docker build -t docker.pkg.github.com/labd/mach-composer/mach:latest --build-arg PROVIDER=full . --progress=plain
