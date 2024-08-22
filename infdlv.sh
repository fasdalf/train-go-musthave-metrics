#!/bin/bash
go install github.com/go-delve/delve/cmd/dlv@latest; go mod tidy; go mod vendor

while true
do
	echo "Press [CTRL+C] to stop..."
	sleep 1
	dlv debug ./cmd/server/main.go --headless=true --api-version=2  -- -k 1234 --filestoragepath new.json --restore
done