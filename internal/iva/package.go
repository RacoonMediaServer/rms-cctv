package iva

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-cctv/internal/model"
)

type PackedEvent struct {
	content interface{}
}

func PackEvent[T Event | Malfunction](content *T) *PackedEvent {
	return &PackedEvent{content: content}
}

func (e *PackedEvent) SetCameraId(id model.CameraID) {
	switch event := e.content.(type) {
	case *Event:
		event.CameraId = id
	case *Malfunction:
		event.CameraId = id
	default:
		panic("invalid event")
	}
}

func (e *PackedEvent) CameraId() model.CameraID {
	switch event := e.content.(type) {
	case *Event:
		return event.CameraId
	case *Malfunction:
		return event.CameraId
	default:
		panic("invalid event")
	}
}

func (e *PackedEvent) IsMalfunction() bool {
	_, ok := e.content.(*Malfunction)
	return ok
}

func (e *PackedEvent) AsMalfunction() *Malfunction {
	return e.content.(*Malfunction)
}

func (e *PackedEvent) IsEvent() bool {
	_, ok := e.content.(*Event)
	return ok
}

func (e *PackedEvent) AsEvent() *Event {
	return e.content.(*Event)
}

func (e *PackedEvent) String() string {
	switch event := e.content.(type) {
	case *Event:
		return fmt.Sprintf("%+v", *event)
	case *Malfunction:
		return fmt.Sprintf("%+v", *event)
	default:
		panic("invalid event")
	}
}
