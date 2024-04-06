package db

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-cctv/internal/service/schedules"
	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	conn *gorm.DB
}

func Connect(config configuration.Database) (*Database, error) {
	db, err := gorm.Open(postgres.Open(config.GetConnectionString()))
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&model.Camera{}, &model.Schedule{}); err != nil {
		return nil, err
	}

	result := Database{conn: db}
	_ = result.AddSchedule(schedules.DefaultSchedule)

	return &result, nil
}
