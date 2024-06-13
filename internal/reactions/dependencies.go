package reactions

import "github.com/RacoonMediaServer/rms-packages/pkg/schedule"

type Registry interface {
	Find(id string, defaultIfNotExists bool) *schedule.Schedule
}
