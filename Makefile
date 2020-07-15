GOOS?=linux
GOARCH?=amd64
ENV?=dev

NAME=mediaserver
VERSION?=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)

REGISTRY_SERVER?=registry.videocoin.net
REGISTRY_PROJECT?=cloud

.PHONY: deploy

default: build

version:
	@echo ${VERSION}

build:
	GOOS=${GOOS} GOARCH=${GOARCH} \
		go build \
			-mod vendor \
			-ldflags="-w -s -X main.Version=${VERSION}" \
			-o bin/${NAME} \
			./cmd/main.go

deps:
	GO111MODULE=on go mod vendor
	cp -r ${GOPATH}/src/github.com/notedit/media-server-go/wrapper ./vendor/github.com/notedit/media-server-go
	cp -r ${GOPATH}/src/github.com/notedit/media-server-go/include ./vendor/github.com/notedit/media-server-go

lint:
	golangci-lint run -v

release: docker-build docker-push

docker-lint:
	docker build -f Dockerfile.lint .

docker-build-base:
	docker build -t ${REGISTRY_SERVER}/${REGISTRY_PROJECT}/mediaserver-base:1.0 -f Dockerfile.base .

docker-push-base:
	docker push ${REGISTRY_SERVER}/${REGISTRY_PROJECT}/mediaserver-base:1.0

docker-build:
	docker build -t ${REGISTRY_SERVER}/${REGISTRY_PROJECT}/${NAME}:${VERSION} -f Dockerfile .

docker-push:
	docker push ${REGISTRY_SERVER}/${REGISTRY_PROJECT}/${NAME}:${VERSION}

deploy:
	cd deploy && helm upgrade -i --wait --set image.tag="${VERSION}" -n console mediaserver ./helm
