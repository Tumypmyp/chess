all: build run

build-dependencies:
	docker build -t dependencies -f ./dependencies.Dockerfile .

build: build-dependencies
	docker-compose -f docker-compose.yaml build

run:
	docker-compose -f docker-compose.yaml up