#!/bin/bash
set -e
CGO_ENABLED=1 GOOS=linux go build -a -o bin/atlas-mapserver ./cmd/mapServer/mapServer.go
CGO_ENABLED=1 GOOS=linux go build -a -o bin/atlas-mockPost ./cmd/mockPost/mockPost.go
docker build -t antihax/atlas-mapserver .
