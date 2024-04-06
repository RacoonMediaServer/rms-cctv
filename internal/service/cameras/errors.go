package cameras

import (
	"fmt"
	"go-micro.dev/v4/logger"
)

func makeError(l logger.Logger, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	l.Log(logger.ErrorLevel, err)
	return err
}
