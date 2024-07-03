#!/bin/bash
echo "testing iteration \"$1\"."
echo "build agent ..."
go build -o agent cmd/agent/main.go && \
echo "build server ..." && \
go build -o server cmd/server/main.go && \
echo "starting metricstest ..." && \
./metricstest -test.v -test.run=^TestIteration$1$ -agent-binary-path=./agent -binary-path=./server