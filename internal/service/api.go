package service

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/teambition/rrule-go"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) GetSettings(ctx context.Context, empty *emptypb.Empty, settings *rms_cctv.CctvSettings) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) SetSettings(ctx context.Context, settings *rms_cctv.CctvSettings, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetCameras(ctx context.Context, empty *emptypb.Empty, response *rms_cctv.GetCamerasResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) AddCamera(ctx context.Context, c *rms_cctv.Camera, response *rms_cctv.AddCameraResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "service",
		"camera": c.Name,
		"method": "AddCamera",
	})

	schedule, err := rrule.StrToRRuleSet(c.Schedule)
	if err != nil {
		return makeError(l, "parse camera schedule failed: %w", err)
	}

	cam := model.Camera{Info: c}
	if err := s.CameraManager.Register(&cam); err != nil {
		return makeError(l, "register camera on the external CCTV system failed: %w", err)
	}

	// TODO: add to database
	// XXX: get camera ID
	const cameraId = 0
	cam.ID = cameraId
	cam.Info.Id = cameraId

	if err := s.CameraManager.Add(&cam, s.Reactor.PushEvent); err != nil {
		// TODO: drop from database
		if err := s.CameraManager.Unregister(&cam); err != nil {
			l.Logf(logger.ErrorLevel, "Unregister camera failed: %s", err)
		}
		return makeError(l, "register camera failed: %w", err)
	}

	s.Reactor.SetReactions(cam.ID, s.makeEventReactions(cam.Info, schedule))

	l.Log(logger.InfoLevel, "Camera added")
	response.CameraId = cameraId
	return nil
}

func (s Service) ModifyCamera(ctx context.Context, c *rms_cctv.Camera, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteCamera(ctx context.Context, request *rms_cctv.DeleteCameraRequest, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetLiveUri(ctx context.Context, request *rms_cctv.GetLiveUriRequest, response *rms_cctv.GetLiveUriResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetReplayUri(ctx context.Context, request *rms_cctv.GetReplayUriRequest, request2 *rms_cctv.GetReplayUriRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetSnapshot(ctx context.Context, request *rms_cctv.GetSnapshotRequest, response *rms_cctv.GetSnapshotResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "service",
		"camera": request.CameraId,
		"method": "GetSnapshot",
	})
	cam, err := s.CameraManager.GetCamera(request.CameraId)
	if err != nil {
		return makeError(l, "operation failed: %w", err)
	}
	response.Snapshot, err = cam.TakeSnapshot(model.PrimaryProfile)
	if err != nil {
		return makeError(l, "get snapshot failed: %w", err)
	}
	return nil
}
