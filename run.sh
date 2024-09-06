#!/usr/bin/env bash
set -e

#### 项目路径：
export ProjectRoot=$(cd "$(dirname "$0")"; pwd)
#### 公共配置地址：
export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
#### agora-go-sdk
export LD_LIBRARY_PATH=$ProjectRoot/pkg/agora-go-sdk/agora_sdk:$LD_LIBRARY_PATH
#### 微软speech-sdk环境变量：
export CGO_CFLAGS="-I$ProjectRoot/pkg/microsoft/speechsdk/include/c_api"
export CGO_LDFLAGS="-L$ProjectRoot/pkg/microsoft/speechsdk/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH=$ProjectRoot/pkg/microsoft/speechsdk/lib/x64:$LD_LIBRARY_PATH

go build -ldflags "-X 'main.buildTimeStamp=$(date +%s)'" -o main.out main.go
./main.out $@
