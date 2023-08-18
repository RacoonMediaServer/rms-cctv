package model

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/cctv"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
)

type Camera struct {
	ID                        uint32
	Info                      *rms_cctv.Camera
	PrimaryProfileToken       string
	PrimaryExternalStreamID   cctv.ID
	SecondaryProfileToken     string
	SecondaryExternalStreamID cctv.ID
	ExternalArchiveID         cctv.ID
}
