package service

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/manager"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactions"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	DeviceManager DeviceManager
	Reactor       Reactor
	Notifier      micro.Publisher
	ReactFactory  reactions.Factory
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

	u, err := parseCameraUrl(c)
	if err != nil {
		return makeError(l, "parse camera URL failed: %w", err)
	}

	// TODO: add to database

	// XXX: get and set camera ID
	const cameraId = 0

	device := manager.Device{
		Id:   cameraId,
		Name: c.Name,
		Url:  u,
	}

	if err := s.DeviceManager.Add(device, s.Reactor.PushEvent); err != nil {
		return makeError(l, "%w", err)
	}

	s.Reactor.SetReactions(cameraId, s.makeEventReactions(c.Mode, c.Schedule))

	logger.Infof("Camera added", c.Name)
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
	//TODO implement me
	panic("implement me")
}
