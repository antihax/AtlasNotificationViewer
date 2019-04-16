#!/bin/bash
set +e
CGO_ENABLED=1 GOOS=linux go build -a -o bin/atlas-mapserver ./cmd/
docker build -t antihax/atlas-mapserver .
