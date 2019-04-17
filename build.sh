#!/bin/bash
set -e
CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o bin/atlas-mapserver ./cmd/mapServer/mapServer.go
CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o bin/atlas-mockPost ./cmd/mockPost/mockPost.go
docker build -t antihax/atlas-mapserver .
