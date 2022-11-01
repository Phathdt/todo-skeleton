package sdkcm

import (
	goservice "github.com/onpointvn/libs/go-sdk"
)

func Recover(sc goservice.ServiceContext) {
	logger := sc.Logger("service")

	if err := recover(); err != nil {
		logger.Error("recover error", err)
	}
}
