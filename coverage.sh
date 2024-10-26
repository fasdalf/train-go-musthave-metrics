#!/bin/bash

go test ./... -coverprofile cover.out && go tool cover -html=cover.out -o cover.out.html && go tool cover -func cover.out > cover.log

# go test -v -covermode=count -coverprofile=coverage.out ./...