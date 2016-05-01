#!/bin/bash

DOCKER_BIN_PATH=$(type -a docker | awk '{print $3}')
if [ ! -x "$DOCKER_BIN_PATH" ]; then
  echo "Unable to find a docker executable, please install Docker"
  exit 1
fi


echo "Clean old binary build"
rm docker/kodokojo-haproxy-marathon || true
# Test and vet to done work here
#echo "Download and install dependencies"
#docker run --rm -v "$PWD"/dep:/usr/go/ -v "$PWD":/usr/go/src/github.com/kodokojo/kodokojo-haproxy-marathon -w /usr/go/src/github.com/kodokojo/kodokojo-haproxy-marathon -e "GOPATH=/usr/go/" golang:1.6 go get ./...
#go get ./...
echo "Running tests"
#docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=windows -e GOARCH=386 golang:1.6 go test github.com/kodokojo/kodokojo-haproxy-marathon/...
#docker run --rm -v "$PWD"/dep:/usr/go/ -v "$PWD":/usr/go/src/github.com/kodokojo/kodokojo-haproxy-marathon -w /usr/go/src/github.com/kodokojo/kodokojo-haproxy-marathon -e GOPATH=/usr/go/ -e GOOS=linux -e GOARCH=amd64 golang:1.6 go test github.com/kodokojo/kodokojo-haproxy-marathon/...
go test ./...

echo "Build linux amd64 binary"
#docker run --rm v "$PWD"/dep:/usr/go/ -v "$PWD":/usr/go/src/github.com/kodokojo/kodokojo-haproxy-marathon -w /usr/go/src/github.com/kodokojo/kodokojo-haproxy-marathon -e GOPATH=/usr/go/ -e GOOS=linux -e GOARCH=amd64 golang:1.6 build -o docker/kodokojo-haproxy-marathon github.com/kodokojo/kodokojo-haproxy-marathon
env GOOS=linux GOARCH=amd64 go build -o docker/kodokojo-haproxy-marathon github.com/kodokojo/kodokojo-haproxy-marathon

echo "Build Docker image"
docker build -t="kodokojo/kodokojo-haproxy-marathon" --no-cache ./docker/
