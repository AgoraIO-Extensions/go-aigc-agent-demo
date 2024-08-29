package workerid

import "github.com/google/uuid"

var UUID = uuid.New().String() // 这只是默认值，实际会被一些初始化操作修改
