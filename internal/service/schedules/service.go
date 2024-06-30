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
	Registry Registry
}

func (s Service) Initialize() error {
	l := logger.Fields(map[string]interface{}{
		"from":   "schedules",
		"method": "Initialize",
	})
	list, err := s.Database.LoadSchedules()
	if err != nil {
		return err
	}
	for _, sched := range list {
		parsed, err := schedule.Parse(sched.Intervals)
		if err != nil {
			l.Logf(logger.ErrorLevel, "Parse schedule '%s:%s' failed: %s", sched.ID, sched.Name, err)
			continue
		}
		s.Registry.Store(sched.ID, parsed)
	}
	return nil
}

func (s Service) GetSchedulesList(ctx context.Context, empty *emptypb.Empty, response *rms_cctv.GetScheduleListResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "schedules",
		"method": "GetSchedulesList",
	})
	list, err := s.Database.LoadSchedules()
	if err != nil {
		l.Logf(logger.ErrorLevel, "Load schedules failed: %s", err)
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
	l.Log(logger.DebugLevel, "Request")
	return nil
}

func (s Service) AddSchedule(ctx context.Context, reqSchedule *rms_cctv.Schedule, response *rms_cctv.AddScheduleResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "schedules",
		"name":   reqSchedule.Name,
		"method": "AddShedule",
	})
	id := uuid.NewString()
	parsed, err := schedule.Parse(reqSchedule.Content)
	if err != nil {
		l.Logf(logger.ErrorLevel, "Parse schedule failed: %s", err)
		return err
	}
	sched := model.Schedule{
		ID:        id,
		Name:      reqSchedule.Name,
		Intervals: reqSchedule.Content,
		Schedule:  parsed,
	}
	if err = s.Database.AddSchedule(&sched); err != nil {
		l.Logf(logger.ErrorLevel, "Add schedule to database failed: %s", err)
		return err
	}
	s.Registry.Store(id, parsed)

	response.Id = id
	l.Logf(logger.InfoLevel, "%s added, %+v", id, sched.Intervals)
	return nil
}

func (s Service) DeleteSchedule(ctx context.Context, request *rms_cctv.DeleteScheduleRequest, empty *emptypb.Empty) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "schedules",
		"id":     request.Id,
		"method": "DeleteSchedule",
	})
	if request.Id == "default" {
		return errors.New("cannot remove default schedule")
	}

	err := s.Database.RemoveSchedule(request.Id)
	if err == nil {
		l.Log(logger.InfoLevel, "Removed")
		s.Registry.Delete(request.Id)
	} else {
		l.Logf(logger.WarnLevel, "Remove failed: %s", err)
	}
	return err
}

func (s Service) GetSchedule(ctx context.Context, request *rms_cctv.GetScheduleRequest, schedule *rms_cctv.Schedule) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "schedules",
		"id":     request.Id,
		"method": "GetSchedules",
	})
	result, err := s.Database.GetSchedule(request.Id)
	if err != nil {
		l.Logf(logger.ErrorLevel, "Fetch schedule failed: %s", err)
		return err
	}

	*schedule = rms_cctv.Schedule{
		Id:      result.ID,
		Name:    result.Name,
		Content: result.Intervals,
	}

	l.Log(logger.DebugLevel)

	return nil
}

func (s Service) ModifySchedule(ctx context.Context, request *rms_cctv.Schedule, empty *emptypb.Empty) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "schedules",
		"id":     request.Id,
		"method": "ModifySchedule",
	})
	parsed, err := schedule.Parse(request.Content)
	if err != nil {
		l.Logf(logger.ErrorLevel, "Parse schedule failed: %s", err)
		return err
	}
	sched := model.Schedule{
		ID:        request.Id,
		Name:      request.Name,
		Intervals: request.Content,
		Schedule:  parsed,
	}
	if err = s.Database.UpdateSchedule(&sched); err != nil {
		l.Logf(logger.ErrorLevel, "Add schedule to database failed: %s", err)
		return err
	}
	s.Registry.Store(request.Id, parsed)

	l.Logf(logger.InfoLevel, "Modified %+v", sched)
	return nil
}

func (s Service) FindSchedule(id string, defaultIfNotExists bool) *model.Schedule {
	result, err := s.Database.GetSchedule(id)
	if err != nil {
		if defaultIfNotExists {
			return DefaultSchedule
		}
		return nil
	}
	return result
}
