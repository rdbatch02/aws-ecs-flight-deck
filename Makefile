.ONESHELL:
build: build-refresh build-read-api

build-refresh:
	cd src/refresh-lambda && \
	rm -rf build && \
	GOOS=linux go build -o build/refresh refresh.go && \
	cd build && zip refresh.zip build/refresh

build-read-api:
	cd src/cluster-read-api && \
	rm -rf build && \
	GOOS=linux go build -o build/ReadApi ReadApi.go && \
	cd build && zip readapi.zip ReadApi

deploy:
	cdk deploy