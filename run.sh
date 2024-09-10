#!/usr/bin/env bash
set -e

#### Project root path：
export ProjectRoot=$(cd "$(dirname "$0")"; pwd)
#### Path to system library files：
export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
#### agora-go-sdk
export LD_LIBRARY_PATH=$ProjectRoot/pkg/agora-go-sdk/agora_sdk:$LD_LIBRARY_PATH
#### Microsoft Speech SDK environment variables：
export CGO_CFLAGS="-I$ProjectRoot/pkg/microsoft/speechsdk/include/c_api"
export CGO_LDFLAGS="-L$ProjectRoot/pkg/microsoft/speechsdk/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH=$ProjectRoot/pkg/microsoft/speechsdk/lib/x64:$LD_LIBRARY_PATH

go build -ldflags "-X 'main.buildTimeStamp=$(date +%s)'" -o main.out main.go
./main.out $@
