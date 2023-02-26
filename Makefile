#!/usr/bin/env make

.PHONY: build
build:
	@ mkdir -p bin
	go build -o bin/rate-limit-server cmd/rate-limit-server/main.go

.PHONY: run
run: build
	./bin/rate-limit-server

.PHONY: docker-build
docker-build:
	docker build -t=bhargav0infracloudio/rate-limit-server:latest .

.PHONY: docker-run
docker-run: docker-build
	docker run --rm --name rate-limit-server -p 8080:8080 bhargav0infracloudio/rate-limit-server:latest

.PHONY: docker-push
docker-push: docker-build
	docker push bhargav0infracloudio/rate-limit-server:latest