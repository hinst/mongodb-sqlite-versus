#!/bin/bash
export CGO_ENABLED=1
export CGO_LDFLAGS="-ldl"
go build && ./mongodb-sqlite-versus "$@"