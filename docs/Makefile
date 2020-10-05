.PHONY: all

install:
	pip install -r requirements.txt

preview:
	mkdocs serve

docker-image:
	docker build -t docs .

docker-build:
	docker run --name docs docs mkdocs build -d /build
	docker cp src:/build .
	docker rm -f docs

clean:
	@rm -rf build
