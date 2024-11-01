#!/bin/bash
go vet -vettool=statictest ./... && \
go build -o=cmd/staticlint/staticlint ./cmd/staticlint && \
go vet -vettool=cmd/staticlint/staticlint ./...