.ONESHELL:
build: build-refresh build-read-api build-refresh-clock build-ui

build-refresh:
	cd src/refresh-lambda && \
	rm -rf build && \
	GOOS=linux go build -o build/refresh refresh.go && \
	cd build && zip refresh.zip refresh

build-read-api:
	cd src/cluster-read-api && \
	rm -rf build && \
	GOOS=linux go build -o build/ReadApi ReadApi.go && \
	cd build && zip readapi.zip ReadApi

build-refresh-clock:
	cd src/refresh-clock && \
	rm -rf build && \
	GOOS=linux go build -o build/RefreshClock RefreshClock.go && \
	cd build && zip refreshclock.zip RefreshClock

build-ui:
	cd src/ui && \
	npm run build

deploy:
	cdk deploy