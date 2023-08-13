package reactions

import (
	"context"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"time"
)

const notifyTimeout = 5 * time.Second

func notify(l logger.Logger, pub micro.Publisher, notification interface{}) {
	go func() {
		ctx, cancel := context.WithTimeout(context.TODO(), notifyTimeout)
		defer cancel()

		if err := pub.Publish(ctx, notification); err != nil {
			l.Logf(logger.ErrorLevel, "Notification failed: %s", err)
			return
		}
		l.Log(logger.InfoLevel, "Notification sent")
	}()
}
