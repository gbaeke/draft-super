# Makefile for go-template

REPO:=gbaeke
TAG:=latest
IMAGE:=$(REPO)/super:$(TAG)


test:
	go test -v -race ./...

build:
	CGO_ENABLED=0 go build -installsuffix 'static' -o app cmd/app/*

docker-build:
	docker build -t $(IMAGE) .

docker-push:
	docker build -t $(IMAGE) .
	docker push $(IMAGE)

swagger:
	cd pkg/api ; swag init -g server.go

dapr-mqtt:
	dapr run --dapr-http-port 3500 --app-id goapp --app-port 8080 -d ./components ./app

dapr:
	dapr run --dapr-http-port 3500 --app-id goapp --app-port 8080 ./app

dapr2:
	PORT=8081 dapr run --dapr-http-port 3501 --app-id goapp2 --app-port 8081 ./app