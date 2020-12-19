.ONESHELL:
build: build-refresh

build-refresh:
	cd src/refresh-lambda && \
	GOOS=linux go build refresh.go && \
	rm -rf refresh.zip && \
	zip refresh.zip refresh

deploy:
	cdk deploy