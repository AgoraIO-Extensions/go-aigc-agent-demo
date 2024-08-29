#!/usr/bin/env bash
set -e

#### 项目路径：
export ProjectRoot=$(cd "$(dirname "$0")"; pwd)
#### 公共配置地址：
export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
#### agora-go-sdk
export LD_LIBRARY_PATH=$ProjectRoot/pkg/agora-go-sdk/agora_sdk:$LD_LIBRARY_PATH

go build -ldflags "-X 'main.buildTimeStamp=$(date +%s)'" -o main.out main.go
./main.out $@
