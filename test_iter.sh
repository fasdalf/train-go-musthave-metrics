#!/bin/bash
echo "testing iteration $1"
 go build -o agent cmd/agent/main.go && \
 go build -o server cmd/server/main.go && \
 ./metricstest -test.v -test.run=^TestIteration$1$ -agent-binary-path=./agent -binary-path=./server