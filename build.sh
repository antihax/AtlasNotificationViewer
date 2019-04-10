#!/bin/bash
GOOS=windows GOARCH=386 go build -o ./dist/hello.exe hello.go
GOOS=linux  go build -o ./dist/hello hello.go