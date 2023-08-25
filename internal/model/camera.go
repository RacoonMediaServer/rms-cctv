package model

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/cctv"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
)

type CameraID uint32

type Camera struct {
	ID                        CameraID         `gorm:"primaryKey;autoIncrement"`
	Info                      *rms_cctv.Camera `gorm:"embedded"`
	PrimaryProfileToken       string
	PrimaryExternalStreamID   cctv.ID
	SecondaryProfileToken     string
	SecondaryExternalStreamID cctv.ID
	ExternalArchiveID         cctv.ID
}
