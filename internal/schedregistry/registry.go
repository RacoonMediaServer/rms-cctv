package schedregistry

import (
	"sync"

	"github.com/RacoonMediaServer/rms-packages/pkg/schedule"
)

type Registry struct {
	storage sync.Map
}

func (r *Registry) Store(id string, sched *schedule.Schedule) {
	r.storage.Store(id, sched)
}

func (r *Registry) Find(id string, defaultIfNotExists bool) *schedule.Schedule {
	empty, _ := schedule.New(schedule.Representation{})
	result, ok := r.storage.Load(id)
	if !ok {
		if !defaultIfNotExists {
			return nil
		}
		return empty
	}

	sched, ok := result.(*schedule.Schedule)
	if !ok {
		if !defaultIfNotExists {
			return nil
		}
		return empty
	}

	return sched
}

func (r *Registry) Delete(id string) {
	r.storage.Delete(id)
}
