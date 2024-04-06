package schedules

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/schedule"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	Database Database
}

func (s Service) GetSchedulesList(ctx context.Context, empty *emptypb.Empty, response *rms_cctv.GetScheduleListResponse) error {
	list, err := s.Database.LoadSchedules()
	if err != nil {
		logger.Errorf("Load schedules failed: %s", err)
		return err
	}

	response.Result = make([]*rms_cctv.Schedule, len(list))
	for i, sch := range list {
		response.Result[i] = &rms_cctv.Schedule{
			Id:      sch.ID,
			Name:    sch.Name,
			Content: sch.Intervals,
		}
	}
	return nil
}

func (s Service) AddSchedule(ctx context.Context, reqSchedule *rms_cctv.Schedule, response *rms_cctv.AddScheduleResponse) error {
	id := uuid.NewString()
	parsed, err := schedule.Parse(reqSchedule.Content)
	if err != nil {
		logger.Errorf("Parse schedule failed: %s", err)
		return err
	}
	sched := model.Schedule{
		ID:        id,
		Name:      reqSchedule.Name,
		Intervals: reqSchedule.Content,
		Schedule:  parsed,
	}
	if err = s.Database.AddSchedule(&sched); err != nil {
		logger.Errorf("Add schedule to database failed: %s", err)
		return err
	}
	response.Id = id
	return nil
}

func (s Service) DeleteSchedule(ctx context.Context, request *rms_cctv.DeleteScheduleRequest, empty *emptypb.Empty) error {
	if request.Id == "default" {
		return errors.New("cannot remove default schedule")
	}
	return s.Database.RemoveSchedule(request.Id)
}

func (s Service) GetSchedule(id string, defaultIfNotExists bool) *model.Schedule {
	result, err := s.Database.GetSchedule(id)
	if err != nil {
		if defaultIfNotExists {
			return DefaultSchedule
		}
		return nil
	}
	return result
}
