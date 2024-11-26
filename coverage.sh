#!/bin/bash

go test ./... -bench=. -coverprofile cover.tmp.out && \
grep -v ".pb.go:" cover.tmp.out > cover.out && \
go tool cover -html=cover.out -o cover.out.html && \
go tool cover -func cover.out > cover.log

# go test -v -covermode=count -coverprofile=coverage.out ./...