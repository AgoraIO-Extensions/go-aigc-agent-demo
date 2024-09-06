package workerid

import (
	"github.com/google/uuid"
	"strings"
)

var UUID = strings.Join(strings.Split(uuid.New().String(), "-"), "")
