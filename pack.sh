#!/usr/bin/env bash

cd frontend && npm run build || exit
cd .. || exit

go-bindata -pkg migrationstatic -o staticfile/migrationstatic/migration.go migrations/... || exit

go-bindata-assetfs -pkg webstatic -o staticfile/webstatic/bindata.go web web/css/... web/fonts/... || exit

mkdir -p staticfile/webstatic/img/imgclient; mkdir -p staticfile/webstatic/img/imghelp; mkdir -p staticfile/webstatic/img/imgicons
go-bindata-assetfs -pkg imgclient -o staticfile/webstatic/img/imgclient/bindata.go web/img/client/... || exit
go-bindata-assetfs -pkg imghelp -o staticfile/webstatic/img/imghelp/bindata.go web/img/help/... || exit
go-bindata-assetfs -pkg imgicons -o staticfile/webstatic/img/imgicons/bindata.go web/img/icons/... || exit

mkdir -p staticfile/webstatic/webjs
go-bindata-assetfs -pkg webjs -o staticfile/webstatic/webjs/bindata.go web/js/... || exit

env CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" GOOS=linux GOARCH=amd64 go build ./v2ray-helper.go || exit
tar zcf v2ray-helper-1.0.0.tgz v2ray-helper
