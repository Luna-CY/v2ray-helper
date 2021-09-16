#!/usr/bin/env bash

function build_web() {
    cd frontend && npm run build || exit
    cd .. || exit
}

function build_static() {
    go-bindata -pkg migrationstatic -o staticfile/migrationstatic/migration.go migrations/... || exit
    go-bindata-assetfs -pkg webstatic -o staticfile/webstatic/bindata.go web web/css/... web/fonts/... || exit

    mkdir -p staticfile/webstatic/webjs
    go-bindata-assetfs -pkg webjs -o staticfile/webstatic/webjs/bindata.go web/js/... || exit
}

function build_binary_and_package() {
    env CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" GOOS=linux GOARCH=amd64 go build ./v2ray-helper.go || exit
    tar zcf v2ray-helper.tgz v2ray-helper
}

while getopts "a" opt; do
  case $opt in
    a)
      build_web
      build_static
      ;;
    ?)
      echo "打包工具，此工具默认只编译主项目，不对前端进行重新构建与资源打包，如果需要构建前端并打包静态资源请私用-a选项"
      echo ""
      echo "Usage: ./pack.sh [options]"
      echo "  -a 构建前端并打包所有资源"
      exit 1
      ;;
  esac
done

build_binary_and_package

