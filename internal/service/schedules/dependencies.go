package schedules

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/schedule"
)

type Database interface {
	LoadSchedules() ([]*model.Schedule, error)
	AddSchedule(sched *model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	UpdateSchedule(sched *model.Schedule) error
	RemoveSchedule(id string) error
}

type Registry interface {
	Store(id string, sched *schedule.Schedule)
	Delete(id string)
}
