package cctv

import (
	"context"
	"fmt"
	cctv_backend "github.com/RacoonMediaServer/rms-packages/pkg/service/cctv-backend"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/RacoonMediaServer/rms-packages/pkg/video"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/url"
	"time"
)

type externalBackend struct {
	l                logger.Logger
	streamService    cctv_backend.StreamService
	recordingService cctv_backend.RecordingService
	systemService    cctv_backend.SystemService
}

const reqTimeout = 1 * time.Minute

func (b externalBackend) AddStream(camera *rms_cctv.Camera, u *url.URL) (ID, error) {
	req := cctv_backend.AddStreamRequest{
		AdviceId: camera.Name,
		Url:      u.String(),
	}

	resp, err := b.streamService.AddStream(context.TODO(), &req, client.WithRequestTimeout(reqTimeout))
	if err != nil {
		return "", err
	}

	return ID(resp.StreamId), nil
}

func (b externalBackend) DeleteStream(id ID) error {
	_, err := b.streamService.RemoveStream(context.TODO(), &cctv_backend.RemoveStreamRequest{StreamId: string(id)}, client.WithRequestTimeout(reqTimeout))
	return err
}

func (b externalBackend) GetStreamUri(id ID, transport video.Transport) (*url.URL, error) {
	req := cctv_backend.GetStreamUriRequest{
		StreamId:  string(id),
		Transport: &transport,
	}

	resp, err := b.streamService.GetStreamUri(context.TODO(), &req, client.WithRequestTimeout(reqTimeout))
	if err != nil {
		return &url.URL{}, err
	}

	u, err := url.Parse(resp.Uri)
	if err != nil {
		return &url.URL{}, fmt.Errorf("parse URL '%s' failed: %w", resp.Uri, err)
	}

	return u, nil
}

func (b externalBackend) StartRecording(id ID) error {
	req := cctv_backend.SetRecordingStateRequest{
		RecordingId: string(id),
		Pause:       false,
	}
	_, err := b.recordingService.SetRecordingState(context.TODO(), &req, client.WithRequestTimeout(reqTimeout))
	return err
}

func (b externalBackend) StopRecording(id ID) error {
	req := cctv_backend.SetRecordingStateRequest{
		RecordingId: string(id),
		Pause:       true,
	}
	_, err := b.recordingService.SetRecordingState(context.TODO(), &req, client.WithRequestTimeout(reqTimeout))
	return err
}

func (b externalBackend) SetQuality(id ID, quality uint) error {
	req := cctv_backend.SetRecordingQualityRequest{
		RecordingId: string(id),
		Quality:     uint32(quality),
	}
	_, err := b.recordingService.SetRecordingQuality(context.TODO(), &req, client.WithRequestTimeout(reqTimeout))
	return err
}

func (b externalBackend) AddArchive(stream ID, rotationDays uint) (ID, error) {
	req := cctv_backend.AddRecordingRequest{
		StreamId:     string(stream),
		RotationDays: uint32(rotationDays),
	}

	resp, err := b.recordingService.AddRecording(context.TODO(), &req, client.WithRequestTimeout(reqTimeout))
	if err != nil {
		return "", err
	}

	return ID(resp.RecordingId), nil
}

func (b externalBackend) DeleteArchive(id ID) error {
	req := cctv_backend.RemoveRecordingRequest{
		RecordingId:      string(id),
		SaveRecordedData: true,
	}

	_, err := b.recordingService.RemoveRecording(context.TODO(), &req, client.WithRequestTimeout(reqTimeout))
	return err
}

func (b externalBackend) GetReplayUri(id ID, transport video.Transport, timestamp time.Time) (*url.URL, error) {
	req := cctv_backend.GetRecordingUriRequest{
		StreamId:  string(id),
		Transport: &transport,
		Timestamp: timestamppb.New(timestamp),
	}

	resp, err := b.recordingService.GetRecordingUri(context.TODO(), &req, client.WithRequestTimeout(reqTimeout))
	if err != nil {
		return &url.URL{}, err
	}

	u, err := url.Parse(resp.Uri)
	if err != nil {
		return &url.URL{}, fmt.Errorf("parse URL '%s' failed: %w", resp.Uri, err)
	}

	return u, nil
}

func newExternalBackend(l logger.Logger, f servicemgr.ClientFactory) *externalBackend {
	b := externalBackend{l: l}

	sf := servicemgr.NewServiceFactory(f)
	b.streamService = sf.NewCctvStreamService()
	b.recordingService = sf.NewCctvRecordingService()
	b.systemService = sf.NewCctvSystemService()

	return &b
}
