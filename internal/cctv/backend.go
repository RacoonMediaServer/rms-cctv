package cctv

import "net/url"

type ID string

// StreamService is an interface to external CCTV streamer
type StreamService interface {
	// AddStream registers camera on external CCTV system
	AddStream(streamUrl *url.URL) (ID, error)
	// DeleteStream removes camera stream
	DeleteStream(id ID) error
	// GetStreamUri returns stream live URL
	GetStreamUri(id ID) (*url.URL, error)
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

type Backend interface {
	StreamService
	RecorderService
}
