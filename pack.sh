#!/usr/bin/env bash

env CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" GOOS=linux GOARCH=amd64 go build ./v2ray-subscription-server.go | exit
env CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" GOOS=linux GOARCH=amd64 go build ./v2ray-subscription-migrate.go | exit
tar zcf target.tgz ./v2ray-subscription-server ./v2ray-subscription-migrate ./web ./migrations ./config
