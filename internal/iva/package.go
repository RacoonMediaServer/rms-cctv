package iva

type PackedEvent struct {
	content interface{}
}

func PackEvent[T Event | Malfunction](content *T) *PackedEvent {
	return &PackedEvent{content: content}
}

func (e *PackedEvent) SetCameraId(id uint32) {
	switch event := e.content.(type) {
	case *Event:
		event.CameraId = id
	case *Malfunction:
		event.CameraId = id
	default:
		panic("invalid event")
	}
}

func (e *PackedEvent) CameraId() uint32 {
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
