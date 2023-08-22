package settings

import rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"

type Loader interface {
	Load() *rms_cctv.CctvSettings
}

type Saver interface {
	Save(settings *rms_cctv.CctvSettings)
}

type Provider interface {
	Loader
	Saver
}
