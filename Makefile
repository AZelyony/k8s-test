APP=$(shell basename $(shell git remote get-url origin))
#REGISTRY=gcr.io/azelyony
REGISTRY=azelyony
VERSION=$(shell git describe --tags --abbrev=0 --always)-$(shell git rev-parse --short HEAD)
TARGETOS=linux #darwin windows
TARGETARCH=amd64 # arm64

#Use "make build TARGETOS=windows TARGETARCH=amd64"

format: 
	gofmt -s -w ./

get:
	go get

lint:
	golint

test: 
	go test -v

build: format get
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o k8s-test -ldflags "-X=github.com/AZelyony/k8s-test/cmd.appVersion=${VERSION}"

image:
	docker build . -t ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}  --build-arg TARGETARCH=${TARGETARCH}

push:
	docker push ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}

clean: 
	rm -rf k8s-test
	docker rmi ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}