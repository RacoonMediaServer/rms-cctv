package reactions

import (
	"errors"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/mocks"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/golang/mock/gomock"
	"go-micro.dev/v4/logger"
	"testing"
)

func TestErrorReaction_React(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pub := mocks.NewMockEvent(ctrl)
	r := errorReaction{pub: pub}
	r.React(logger.DefaultLogger, *iva.PackEvent(&iva.Event{}))

	pub.EXPECT().Publish(gomock.Any(), gomock.AssignableToTypeOf(&events.Malfunction{}), gomock.Any()).Return(nil).Times(1)
	m := iva.NewMalfunction(errors.New("test"))
	r.React(logger.DefaultLogger, *iva.PackEvent(m))

	pub.EXPECT().Publish(gomock.Any(), gomock.AssignableToTypeOf(&events.Malfunction{}), gomock.Any()).Return(errors.New("failed")).Times(1)
	r.React(logger.DefaultLogger, *iva.PackEvent(m))
}
