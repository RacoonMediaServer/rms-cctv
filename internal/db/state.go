package db

import "github.com/RacoonMediaServer/rms-cctv/internal/model"

type cctvState struct {
	ID    uint        `gorm:"primaryKey"`
	State model.State `gorm:"embedded"`
}

func (d *Database) LoadState() (*model.State, error) {
	var record cctvState
	if err := d.conn.Model(&cctvState{}).FirstOrCreate(&record, &record).Error; err != nil {
		return nil, err
	}
	return &record.State, nil
}

func (d *Database) SaveState(val *model.State) error {
	return d.conn.Save(&cctvState{ID: 1, State: *val}).Error
}
