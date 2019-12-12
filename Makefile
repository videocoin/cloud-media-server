GOOS?=linux
GOARCH?=amd64

GCP_PROJECT=videocoin-network
ENV?=dev

NAME=mediaserver
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)

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
	cp -r $(GOPATH)/src/github.com/notedit/media-server-go/wrapper \
	vendor/github.com/notedit/media-server-go/wrapper
	cp -r $(GOPATH)/src/github.com/notedit/media-server-go/include \
	vendor/github.com/notedit/media-server-go/include

docker-build-base:
	docker build -t gcr.io/${GCP_PROJECT}/mediaserver-go:${VERSION} -f Dockerfile.base .

docker-push-base:
	docker push gcr.io/${GCP_PROJECT}/mediaserver-go:${VERSION}


docker-build:
	docker build -t gcr.io/${GCP_PROJECT}/${NAME}:${VERSION} -f Dockerfile .

docker-push:
	docker push gcr.io/${GCP_PROJECT}/${NAME}:${VERSION}

release: docker-build docker-push

deploy:
	ENV=${ENV} deploy/deploy.sh