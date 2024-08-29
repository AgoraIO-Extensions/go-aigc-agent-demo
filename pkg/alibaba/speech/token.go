package speech

import (
	"errors"
	"fmt"
	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
)

var TOKEN string

func InitToken(akid, akkey string) error {
	tokenMsg, err := nls.GetToken(nls.DEFAULT_DISTRIBUTE, nls.DEFAULT_DOMAIN, akid, akkey, nls.DEFAULT_VERSION)
	if err != nil {
		return err
	}
	if tokenMsg.TokenResult.Id == "" {
		str := fmt.Sprintf("obtain empty token err:%s", tokenMsg.ErrMsg)
		return errors.New(str)
	}
	TOKEN = tokenMsg.TokenResult.Id
	logger.Inst().Info("ali speech token", zap.String("token", TOKEN))
	return nil
}
