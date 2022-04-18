
hello:
	echo 'Hello, welcome to start a new project.'

build:
	go build ./cmd/project_layout

run:
	./project_layout server

build_and_run:build run

docker_build:
	docker build -f data/deploy/Dockerfile -t project_layout .

compose_run:
	cd data/deploy && \
	docker-compose up -d

compose_down:
	cd data/deploy && \
	docker-compose down
