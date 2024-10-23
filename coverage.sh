#!/bin/bash

go test ./... -coverprofile cover.out && go tool cover -html=cover.out -o cover.html.out && go tool cover -func cover.out

