package schedules

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/schedule"
)

func init() {
	DefaultSchedule.Schedule, _ = schedule.Parse("{}")
}

var DefaultSchedule = &model.Schedule{
	ID:        "default",
	Name:      "Все время",
	Intervals: "{}",
	Schedule:  nil,
}
