#!/bin/bash
# go tool pprof -http=":9090" -seconds=30 http://localhost:8093/debug/pprof/goroutine
# можно использовать опции focus и ignore —
# go tool pprof -focus=myproject -ignore="^runtime|^net/http|^[github.com|^golang.org](http://github.com%7C%5Egolang.org)" base.pprof
# Это позволит сфокусироваться на вызовах, связанных именно с нашим проектом.

# curl -s http://localhost:8093/debug/pprof/heap?seconds=30 > server.heap.result.pprof
# go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

go test ./... -coverprofile cover.out && go tool cover -func cover.out

