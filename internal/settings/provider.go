package settings

import (
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"sync"
)

type provider struct {
	mu       sync.RWMutex
	settings *rms_cctv.CctvSettings
}

func New() Provider {
	return &provider{
		settings: &rms_cctv.CctvSettings{
			EventNotifyThresholdIntervalSec: eventNotifyThresholdIntervalSec,
			OneEventDefaultDurationSec:      oneEventDefaultDurationSec,
		},
	}
}

func (p *provider) Load() *rms_cctv.CctvSettings {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return &rms_cctv.CctvSettings{
		EventNotifyThresholdIntervalSec: p.settings.EventNotifyThresholdIntervalSec,
		OneEventDefaultDurationSec:      p.settings.OneEventDefaultDurationSec,
	}
}

func (p *provider) Save(settings *rms_cctv.CctvSettings) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.settings = &rms_cctv.CctvSettings{
		EventNotifyThresholdIntervalSec: settings.EventNotifyThresholdIntervalSec,
		OneEventDefaultDurationSec:      settings.OneEventDefaultDurationSec,
	}
}
