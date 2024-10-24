#!/bin/bash
 go install golang.org/x/tools/cmd/goimports@latest
goimports -w -v internal
goimports -w -v cmd