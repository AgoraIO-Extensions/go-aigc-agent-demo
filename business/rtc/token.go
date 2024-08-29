package rtc

import (
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"strconv"
)

func GetPublisherToken(appid, appCertificate, channelName string, userId int) (string, error) {
	if appCertificate == "" {
		return "", nil
	}

	tokenExpirationInSeconds := uint32(86400)
	privilegeExpirationInSeconds := uint32(86400)
	return rtctokenbuilder.BuildTokenWithUserAccount(appid, appCertificate, channelName, strconv.Itoa(userId),
		rtctokenbuilder.RolePublisher, tokenExpirationInSeconds, privilegeExpirationInSeconds)
}
