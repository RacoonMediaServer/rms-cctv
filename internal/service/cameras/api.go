package cameras

import (
	"context"
	"time"

	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetCameras(ctx context.Context, empty *emptypb.Empty, response *rms_cctv.GetCamerasResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "cameras",
		"method": "GetCameras",
	})
	l.Logf(logger.DebugLevel, "Request")

	cameras, err := s.Database.LoadCameras()
	if err != nil {
		return makeError(l, "fetch cameras failed: %w", err)
	}
	response.Cameras = make([]*rms_cctv.Camera, len(cameras))
	for i := range cameras {
		response.Cameras[i] = cameras[i].Info
	}
	return nil
}

func (s Service) AddCamera(ctx context.Context, c *rms_cctv.Camera, response *rms_cctv.AddCameraResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "cameras",
		"camera": c.Name,
		"method": "AddCamera",
	})
	l.Logf(logger.DebugLevel, "Request")

	schedule := s.ScheduleRegistry.Find(c.Schedule, false)
	if schedule == nil {
		return makeError(l, "cannot find associating schedule")
	}

	// TODO: очень много вариантов нарушить консистентность, но пока так
	cam := model.Camera{Info: c}
	if err := s.Database.AddCamera(&cam); err != nil {
		return makeError(l, "add camera to database failed: %s", err)
	}
	cam.Info.Id = uint32(cam.ID)

	if err := s.CameraManager.Register(&cam); err != nil {
		if err := s.Database.RemoveCamera(cam.ID); err != nil {
			l.Logf(logger.ErrorLevel, "Remove camera failed: %s", err)
		}
		return makeError(l, "register camera on the external CCTV system failed: %w", err)
	}

	if err := s.Database.UpdateCamera(&cam); err != nil {
		if err := s.CameraManager.Unregister(&cam); err != nil {
			l.Logf(logger.ErrorLevel, "Unregister camera failed: %s", err)
		}
		if err := s.Database.RemoveCamera(cam.ID); err != nil {
			l.Logf(logger.ErrorLevel, "Remove camera failed: %s", err)
		}
		return makeError(l, "add camera to database failed: %w", err)
	}

	if err := s.registerCamera(&cam); err != nil {
		if err := s.Database.RemoveCamera(cam.ID); err != nil {
			l.Logf(logger.ErrorLevel, "Remove camera failed: %s", err)
		}
		if err := s.CameraManager.Unregister(&cam); err != nil {
			l.Logf(logger.ErrorLevel, "Unregister camera failed: %s", err)
		}
		return makeError(l, "initialize camera failed: %w", err)
	}

	l.Logf(logger.InfoLevel, "Camera %s [ %d ] added", cam.Info.Name, cam.ID)
	response.CameraId = uint32(cam.ID)
	return nil
}

func (s Service) ModifyCamera(ctx context.Context, c *rms_cctv.ModifyCameraRequest, empty *emptypb.Empty) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "cameras",
		"camera": c.Id,
		"method": "ModifyCamera",
	})
	l.Logf(logger.DebugLevel, "Request")
	id := model.CameraID(c.Id)

	cam, err := s.Database.GetCamera(id)
	if err != nil {
		return makeError(l, "fetch camera failed: %w", err)
	}

	schedule := s.ScheduleRegistry.Find(c.Schedule, false)
	if schedule == nil {
		return makeError(l, "cannot find associating schedule")
	}

	cam.Info.Name = c.Name
	cam.Info.KeepDays = c.KeepDays
	cam.Info.Mode = c.Mode
	cam.Info.Schedule = c.Schedule

	if err = s.Database.UpdateCamera(cam); err != nil {
		return makeError(l, "update camera in database failed: %w", err)
	}

	if err = s.CameraManager.Modify(id, c.Name, c.KeepDays, c.Mode); err != nil {
		l.Logf(logger.WarnLevel, "modify camera in runtime failed: %s", err)
	}

	s.Reactor.DropReactions(id)
	s.Reactor.SetReactions(id, s.makeEventReactions(cam.Info))

	return nil
}

func (s Service) DeleteCamera(ctx context.Context, request *rms_cctv.DeleteCameraRequest, empty *emptypb.Empty) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "cameras",
		"camera": request.CameraId,
		"method": "DeleteCamera",
	})
	l.Logf(logger.DebugLevel, "Request")

	id := model.CameraID(request.CameraId)
	if err := s.Database.RemoveCamera(id); err != nil {
		return makeError(l, "cannot remove camera from db: %w", err)
	}
	if err := s.CameraManager.Remove(id); err != nil {
		l.Logf(logger.ErrorLevel, "stop camera %d failed: %s", id, err)
	}
	return nil
}

func (s Service) GetLiveUri(ctx context.Context, request *rms_cctv.GetLiveUriRequest, response *rms_cctv.GetUriResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "cameras",
		"camera": request.CameraId,
		"method": "GetLiveUri",
	})
	l.Logf(logger.DebugLevel, "Request")

	profile := model.PrimaryProfile
	if !request.MainProfile {
		profile = model.SecondaryProfile
	}

	uri, err := s.CameraManager.GetStreamUri(model.CameraID(request.CameraId), profile, request.Transport)
	if err != nil {
		return makeError(l, "operation failed: %w", err)
	}

	response.Uri = uri
	return nil
}

func (s Service) GetReplayUri(ctx context.Context, request *rms_cctv.GetReplayUriRequest, response *rms_cctv.GetUriResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "cameras",
		"camera": request.CameraId,
		"method": "GetReplayUri",
	})
	l.Logf(logger.DebugLevel, "Request")

	ts := time.Time{}
	if request.Timestamp != nil {
		ts = time.Unix(int64(*request.Timestamp), 0)
	}

	uri, err := s.CameraManager.GetReplayUri(model.CameraID(request.CameraId), request.Transport, ts)
	if err != nil {
		return makeError(l, "operation failed: %w", err)
	}

	response.Uri = uri
	return nil
}

func (s Service) GetSnapshot(ctx context.Context, request *rms_cctv.GetSnapshotRequest, response *rms_cctv.GetSnapshotResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "cameras",
		"camera": request.CameraId,
		"method": "GetSnapshot",
	})
	l.Logf(logger.DebugLevel, "Request")

	cam, err := s.CameraManager.GetCamera(model.CameraID(request.CameraId))
	if err != nil {
		return makeError(l, "operation failed: %w", err)
	}
	response.Snapshot, err = cam.TakeSnapshot(model.PrimaryProfile)
	if err != nil {
		return makeError(l, "get snapshot failed: %w", err)
	}
	return nil
}
