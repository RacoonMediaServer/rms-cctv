package model

import "github.com/RacoonMediaServer/rms-packages/pkg/schedule"

type Schedule struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	Intervals string
	Schedule  *schedule.Schedule `gorm:"-"`
}
