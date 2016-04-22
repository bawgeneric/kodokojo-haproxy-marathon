#!/bin/bash

echo "Clean old binary build"
rm docker/kodokojo-haproxy-marathon || true
# Test and vet to done work here
echo "Running tests"
go test github.com/kodokojo/kodokojo-haproxy-marathon/... 

echo "Build linux amd64 binary"
env GOOS=linux GOARCH=amd64 go build -o docker/kodokojo-haproxy-marathon github.com/kodokojo/kodokojo-haproxy-marathon

echo "Build Docker image"
docker build -t="kodokojo/kodokojo-haproxy-marathon" --no-cache ./docker/
