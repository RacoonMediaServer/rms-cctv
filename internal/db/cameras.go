package db

import "github.com/RacoonMediaServer/rms-cctv/internal/model"

func (d *Database) LoadCameras() ([]*model.Camera, error) {
	var result []*model.Camera
	if err := d.conn.Find(&result).Error; err != nil {
		return nil, err
	}
	for _, cam := range result {
		cam.Info.Id = uint32(cam.ID)
	}
	return result, nil
}

func (d *Database) AddCamera(camera *model.Camera) error {
	return d.conn.Create(camera).Error
}

func (d *Database) UpdateCamera(camera *model.Camera) error {
	return d.conn.Save(camera).Error
}

func (d *Database) RemoveCamera(id model.CameraID) error {
	return d.conn.Model(&model.Camera{}).Unscoped().Delete(&model.Camera{}, id).Error
}
