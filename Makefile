APP_NAME=spiffe-demo-app
DOCKER_IMAGE=spiffe-demo-app
VERSION=latest


.PHONY: all build image clean run

all: build image

build:
	go build -o ./bin/$(APP_NAME)

image:
	ko build --preserve-import-paths --local .

publish_poc:
	KO_DOCKER_REPO=771189981606.dkr.ecr.us-west-2.amazonaws.com ko build --sbom=none --base-import-paths .

resolve_poc:
	KO_DOCKER_REPO=771189981606.dkr.ecr.us-west-2.amazonaws.com ko resolve --sbom=none --base-import-paths -f k8s/spire/deployment.yaml

run:
	KO_DATA_PATH=./kodata SPIFFE_ENDPOINT_SOCKET=unix:///tmp/spirl/spiffe.sock go run main.go

clean:
	rm -rf ./bin
