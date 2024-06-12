package schedules

import "github.com/RacoonMediaServer/rms-cctv/internal/model"

type Database interface {
	LoadSchedules() ([]*model.Schedule, error)
	AddSchedule(sched *model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	UpdateSchedule(sched *model.Schedule) error
	RemoveSchedule(id string) error
}
