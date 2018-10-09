.PHONY: all test clean docker

default: build

build: test cover
	go build -i -o bin/app

up: docker
	docker-compose up

down:
	docker-compose down

db:
	docker-compose up db

docker:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o bin/appdocker
	docker build -t app .

run: docker
	docker run --rm -p 8080:8080 --network app_cloud-builder app

test:
	go test ./...

cover:
	go test ./... -cover

clean:
	rm -rf bin