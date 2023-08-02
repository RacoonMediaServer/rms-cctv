package camera

import (
	"context"
	"net/url"
)

type Type int

const (
	AutoDetect Type = iota
	Onvif
)

func New(ctx context.Context, u *url.URL, t Type) Controller {
	switch t {
	case AutoDetect:
		fallthrough
	case Onvif:
		return newOnvifController(ctx, u)
	}
	panic("unknown type")
}
