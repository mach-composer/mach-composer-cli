PKGS := $(shell find . -type f -name main.go -maxdepth 3 | xargs -I {} dirname {} | cut -c 7- | awk '{print "github.com/labd/mach-composer/cmd/"$$1}')
GO := go
VERSION ?= $(shell git describe --tags --first-parent --abbrev=0 | cut -c 2-)
GOFLAGS ?= -mod=readonly -ldflags "-s -w -X 'main.version=$(VERSION)-dev' -extldflags '-static'"

check: lint cover

build: tidy
	CGO_ENABLED=0 $(GO) build -a -trimpath -tags netgo $(GOFLAGS) -o bin/ ./cmd/...

clean:
	@$(GO) clean

tidy:
	@$(GO) mod tidy -v

test: tidy
	$(GO) test -race ./...

cover: tidy
	$(GO) test -race -coverprofile=coverage.out -covermode=atomic ./...

cover-html: cover
	$(GO) tool cover -html=coverage.out -o coverage.html

lint: vet gofmt misspell staticcheck ineffassign

vet: | test
	$(GO) vet ./...

staticcheck:
	@$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -checks all,-ST1000 ./...

misspell:
	@$(GO) install github.com/client9/misspell/cmd/misspell@latest
	misspell \
		-locale GB \
		-error \
		*.go

ineffassign:
	@$(GO) install github.com/gordonklaus/ineffassign@latest
	ineffassign ./...

pedantic: check errcheck

errcheck:
	@$(GO) install github.com/kisielk/errcheck@latest
	errcheck -ignoretests ./...

gofmt:
	@mkdir -p output
	@rm -f output/lint.log

	gofmt -d -s . 2>&1 | tee output/lint.log

	@[ ! -s output/lint.log ]

	@rm -fr output

docker:
	docker build -t docker.pkg.github.com/labd/mach-composer/mach:latest . --progress=plain
