package service

import (
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"net/url"
)

func parseCameraUrl(c *rms_cctv.Camera) (*url.URL, error) {
	u, err := url.Parse(c.Url)
	if err != nil {
		return nil, err
	}
	if c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	} else if c.User != "" {
		u.User = url.User(c.User)
	}
	return u, nil
}
