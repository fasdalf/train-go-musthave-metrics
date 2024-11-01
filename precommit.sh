#!/bin/bash
./goimports.sh &&\
./vet.sh && \
./coverage.sh