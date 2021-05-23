#!/bin/sh

# 驱动名称
protocol_name="device-service-template"

build() {
    v=$1
    if [ "$1" = '' ]; then
      echo "缺少版本号 如 v0.0.1-dev, 默认版本号为 latest"
      v='latest'
    fi

    p=$(pwd)
    cd "$p"
    docker build -t "$protocol_name":"$v" -f Dockerfile .
    if [ $? != 0 ];then
       exit 1
    fi
}

build "$1";