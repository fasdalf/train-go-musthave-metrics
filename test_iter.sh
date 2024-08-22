#!/bin/bash
echo "testing iteration \"$1\"."
echo "build agent ..."
go build -o cmd/agent/agent cmd/agent/main.go && \
echo "build server ..." && \
go build -o cmd/server/server cmd/server/main.go && \
echo "starting metricstest ..." && \
./metricstest -test.v -test.run=^TestIteration$1$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -source-path=. \
    -server-port=46011 \
    -file-storage-path=other.json \
    -database-dsn='postgres://postgres:postgresP@$@train-go-musthave-metrics_db:5432/postgres_metrics?sslmode=disable' \
    -key='=${TEMP_FILE}=' \

