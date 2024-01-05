package cctv

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-cctv/internal/cctv/client"
	"github.com/RacoonMediaServer/rms-cctv/internal/cctv/client/recorder_service"
	"github.com/RacoonMediaServer/rms-cctv/internal/cctv/client/stream_service"
	"github.com/RacoonMediaServer/rms-cctv/internal/config"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"go-micro.dev/v4/logger"
	"net/url"
	"time"
)

//go:generate swagger generate client --target ./ --name Cctv --spec ../../api/cctv-backend.yaml --principal models.Principal

type externalBackend struct {
	l    logger.Logger
	auth runtime.ClientAuthInfoWriter
	cli  *client.Cctv
}

func newExternalBackend(l logger.Logger, conf config.Backend) Backend {
	tr := httptransport.New(conf.Host, conf.Path, client.DefaultSchemes)

	return &externalBackend{
		l:    l,
		auth: httptransport.APIKeyAuth("X-Token", "header", conf.Token),
		cli:  client.New(tr, strfmt.Default),
	}
}

func (e externalBackend) AddStream(streamUrl *url.URL) (ID, error) {
	u := streamUrl.String()
	req := stream_service.AddStreamParams{
		Stream:  stream_service.AddStreamBody{URL: &u},
		Context: context.TODO(),
	}

	resp, err := e.cli.StreamService.AddStream(&req, e.auth)
	if err != nil {
		return "", err
	}

	id := ID(*resp.GetPayload().ID)
	e.l.Logf(logger.InfoLevel, "Stream %s [ %s ] registered", id, u)

	return id, nil
}

func (e externalBackend) DeleteStream(id ID) error {
	req := stream_service.DeleteStreamParams{
		ID:      string(id),
		Context: context.TODO(),
	}

	_, err := e.cli.StreamService.DeleteStream(&req, e.auth)
	if err != nil {
		return err
	}

	e.l.Logf(logger.InfoLevel, "Stream %s deleted", id)
	return nil
}

func (e externalBackend) GetStreamUri(id ID, transport rms_cctv.VideoTransport) (*url.URL, error) {
	transportString := transport.String()
	req := stream_service.GetStreamURIParams{
		ID:        string(id),
		Transport: &transportString,
		Context:   context.TODO(),
	}

	resp, err := e.cli.StreamService.GetStreamURI(&req, e.auth)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(resp.Payload.URL)
	if err != nil {
		return nil, fmt.Errorf("parse responded URL '%s' failed: %w", resp.Payload.URL, err)
	}

	e.l.Logf(logger.DebugLevel, "Live URL: %s", resp.Payload.URL)
	return u, nil
}

func (e externalBackend) StartRecording(id ID) error {
	req := recorder_service.RecordingControlParams{
		ID:      string(id),
		Pause:   false,
		Context: context.TODO(),
	}

	_, err := e.cli.RecorderService.RecordingControl(&req, e.auth)
	if err != nil {
		return err
	}

	e.l.Logf(logger.InfoLevel, "Recording %s started", id)
	return nil
}

func (e externalBackend) StopRecording(id ID) error {
	req := recorder_service.RecordingControlParams{
		ID:      string(id),
		Pause:   true,
		Context: context.TODO(),
	}

	_, err := e.cli.RecorderService.RecordingControl(&req, e.auth)
	if err != nil {
		return err
	}

	e.l.Logf(logger.InfoLevel, "Recording %s started", id)
	return nil
}

func (e externalBackend) SetQuality(id ID, quality uint) error {
	req := recorder_service.SetRecordingQualityParams{
		ID:      string(id),
		Value:   int64(quality),
		Context: context.TODO(),
	}

	_, err := e.cli.RecorderService.SetRecordingQuality(&req, e.auth)
	if err != nil {
		return err
	}

	e.l.Logf(logger.InfoLevel, "Recording %s change quality to %d", id, quality)
	return nil
}

func (e externalBackend) AddArchive(stream ID, rotationDays uint) (ID, error) {
	streamID := string(stream)
	rotation := int64(rotationDays)
	req := recorder_service.AddRecordingParams{
		Recording: recorder_service.AddRecordingBody{
			RotationDays: &rotation,
			StreamID:     &streamID,
		},
		Context: context.TODO(),
	}

	resp, err := e.cli.RecorderService.AddRecording(&req, e.auth)
	if err != nil {
		return "", err
	}

	id := ID(*resp.GetPayload().ID)
	e.l.Logf(logger.InfoLevel, "Archive '%s' [ %s ] registered", id, streamID)
	return id, nil
}

func (e externalBackend) DeleteArchive(id ID) error {
	req := recorder_service.DeleteRecordingParams{
		ID:      string(id),
		Context: context.TODO(),
	}
	_, err := e.cli.RecorderService.DeleteRecording(&req, e.auth)
	if err != nil {
		return err
	}

	e.l.Logf(logger.InfoLevel, "Archive '%s' deleted", id)
	return nil
}

func (e externalBackend) GetReplayUri(id ID, transport rms_cctv.VideoTransport, timestamp time.Time) (*url.URL, error) {
	transportString := transport.String()
	ts := timestamp.UTC().Unix()
	req := recorder_service.GetReplayURIParams{
		ID:         string(id),
		Timestamp:  &ts,
		Transport:  &transportString,
		Context:    nil,
		HTTPClient: nil,
	}

	resp, err := e.cli.RecorderService.GetReplayURI(&req, e.auth)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(resp.Payload.URL)
	if err != nil {
		return nil, fmt.Errorf("parse responded URL '%s' failed: %w", resp.Payload.URL, err)
	}

	e.l.Logf(logger.DebugLevel, "Replay URL: %s", resp.Payload.URL)
	return u, nil
}
