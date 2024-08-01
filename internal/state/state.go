package state

import (
	"sync"

	"github.com/RacoonMediaServer/rms-cctv/internal/model"
)

type Storage struct {
	mu    sync.RWMutex
	state model.State
}

func (s *Storage) IsNobodyAtHome() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.state.NobodyAtHome
}

func (s *Storage) Load() model.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.state
}

func (s *Storage) Set(state model.State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = state
}

func (s *Storage) SetNobodyAtHome(active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state.NobodyAtHome = active
}

func (s *Storage) Lock() model.State {
	s.mu.RLock()
	return s.state
}

func (s *Storage) Unlock() {
	s.mu.RUnlock()
}
