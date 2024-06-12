package db

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/schedule"
)

func (d *Database) LoadSchedules() ([]*model.Schedule, error) {
	var result []*model.Schedule
	if err := d.conn.Find(&result).Error; err != nil {
		return nil, err
	}
	for _, sched := range result {
		sched.Schedule, _ = schedule.Parse(sched.Intervals)
	}
	return result, nil
}

func (d *Database) AddSchedule(sched *model.Schedule) error {
	return d.conn.Create(sched).Error
}

func (d *Database) GetSchedule(id string) (*model.Schedule, error) {
	var sched model.Schedule
	if err := d.conn.First(&sched, "id = ?", id).Error; err != nil {
		return nil, err
	}
	sched.Schedule, _ = schedule.Parse(sched.Intervals)
	return &sched, nil
}

func (d *Database) UpdateSchedule(sched *model.Schedule) error {
	return d.conn.Save(sched).Error
}

func (d *Database) RemoveSchedule(id string) error {
	return d.conn.Model(&model.Schedule{}).Unscoped().Where("id = ?", id).Delete(&model.Schedule{}).Error
}
