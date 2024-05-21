package zp

import (
	"go.uber.org/zap"
)

func Err(err error) zap.Field {
	return zap.Any("error", err.Error())
}
