package camera

import (
	"context"
	"net/url"
)

type Type int

type ID uint32

const (
	AutoDetect Type = iota
	Onvif
)

type factory struct {
}

type Factory interface {
	New(ctx context.Context, u *url.URL, t Type) Controller
}

func NewFactory() Factory {
	return &factory{}
}

func (f factory) New(ctx context.Context, u *url.URL, t Type) Controller {
	switch t {
	case AutoDetect:
		fallthrough
	case Onvif:
		return newOnvifController(ctx, u)
	}
	panic("unknown type")
}
