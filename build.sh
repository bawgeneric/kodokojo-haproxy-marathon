#!/bin/bash

rm docker/kodokojo-haproxy-marathon || true
# Test and vet to done work here
env GOOS=linux GOARCH=amd64 go build -o docker/kodokojo-haproxy-marathon

docker build -t="kodokojo/kodokojo-haproxy-marathon" --no-cache ./docker/
