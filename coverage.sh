#!/bin/bash
# go tool pprof -http=":9090" -seconds=30 http://localhost:8093/debug/pprof/goroutine
# curl -s http://localhost:8093/debug/pprof/heap > server.heap.result.pprof
# go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

go test ./... -coverprofile cover.out && go tool cover -func cover.out

