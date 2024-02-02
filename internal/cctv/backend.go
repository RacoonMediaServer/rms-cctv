package cctv

import (
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-packages/pkg/video"
	"net/url"
	"time"
)

type ID string

// StreamService is an interface to external CCTV streamer
type StreamService interface {
	// AddStream registers camera on external CCTV system
	AddStream(camera *rms_cctv.Camera, u *url.URL) (ID, error)
	// DeleteStream removes camera stream
	DeleteStream(id ID) error
	// GetStreamUri returns stream live URL
	GetStreamUri(id ID, transport video.Transport) (*url.URL, error)
}

// RecordingController is an interface for control recording process
type RecordingController interface {
	// StartRecording starts record stream to archive
	StartRecording(id ID) error
	// StopRecording stops recording
	StopRecording(id ID) error
	// SetQuality sets recording quality (0 means normal)
	SetQuality(id ID, quality uint) error
}

// RecorderService is an interface to external CCTV recorder
type RecorderService interface {
	RecordingController
	// AddArchive registers stream for recording
	AddArchive(stream ID, rotationDays uint) (ID, error)
	// DeleteArchive removes camera archive
	DeleteArchive(id ID) error
}

// ReplayService is an interface for playback recordings
type ReplayService interface {
	// GetReplayUri gets URI for replay recording
	GetReplayUri(id ID, transport video.Transport, timestamp time.Time) (*url.URL, error)
}

type Backend interface {
	StreamService
	RecorderService
	ReplayService
}
