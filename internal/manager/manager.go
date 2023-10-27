package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/cctv"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"go-micro.dev/v4/logger"
	"net/url"
	"sync"
	"time"
)

type Manager struct {
	l       logger.Logger
	f       camera.Factory
	backend cctv.Backend

	mu       sync.RWMutex
	channels map[model.CameraID]*channel
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(f camera.Factory, backend cctv.Backend) *Manager {
	m := Manager{
		l:        logger.Fields(map[string]interface{}{"from": "manager"}),
		f:        f,
		backend:  backend,
		channels: map[model.CameraID]*channel{},
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())
	return &m
}

func (m *Manager) Add(cam *model.Camera, consumer camera.EventConsumer) error {
	u, err := parseCameraUrl(cam.Info)
	if err != nil {
		return fmt.Errorf("parse camera url failed: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	dev := m.f.New(m.ctx, u, camera.AutoDetect)
	ch := channel{
		camera:    cam,
		cameraUrl: u,
		l:         camera.NewListener(dev, consumer),
	}
	m.channels[cam.ID] = &ch
	return nil
}

// Register create stream and archive to external CCTV system. Its modify cam content
func (m *Manager) Register(cam *model.Camera) error {
	u, err := parseCameraUrl(cam.Info)
	if err != nil {
		return fmt.Errorf("parse camera url failed: %w", err)
	}

	dev := m.f.New(m.ctx, u, camera.AutoDetect)
	profiles, err := dev.GetProfiles()
	if err != nil {
		return fmt.Errorf("cannot get profiles: %w", err)
	}
	if len(profiles) == 0 {
		return errors.New("camera has not any profile")
	}
	urls := make([]*url.URL, len(profiles))
	for i := range profiles {
		u, err = dev.GetStreamUri(profiles[i])
		if err != nil {
			return fmt.Errorf("cannot get stream URL for profile %s: %w", profiles[i], err)
		}
		urls[i] = u
	}

	cam.PrimaryProfileToken = profiles[0]
	cam.PrimaryExternalStreamID, err = m.backend.AddStream(urls[0])
	if err != nil {
		return fmt.Errorf("register stream for profile %s failed: %w", profiles[0], err)
	}

	if len(profiles) != 1 {
		cam.SecondaryProfileToken = profiles[1]
		cam.SecondaryExternalStreamID, err = m.backend.AddStream(urls[1])
		if err != nil {
			if err := m.Unregister(cam); err != nil {
				m.l.Logf(logger.ErrorLevel, "failed to unregister camera: %s", err)
			}
			return fmt.Errorf("register stream for profile %s failed: %w", profiles[1], err)
		}
	}

	cam.ExternalArchiveID, err = m.backend.AddArchive(cam.PrimaryExternalStreamID, uint(cam.Info.KeepDays))
	if err != nil {
		if err := m.Unregister(cam); err != nil {
			m.l.Logf(logger.ErrorLevel, "failed to unregister camera: %s", err)
		}
		return fmt.Errorf("create archive failed: %w", err)
	}

	return nil
}

func (m *Manager) GetCamera(id model.CameraID) (accessor.Camera, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ch, ok := m.channels[id]
	if !ok {
		return nil, fmt.Errorf("camera %d not found", id)
	}

	a := cameraAccessor{
		cam: m.f.New(m.ctx, ch.cameraUrl, camera.AutoDetect),
		profiles: map[model.Profile]string{
			model.PrimaryProfile:   ch.camera.PrimaryProfileToken,
			model.SecondaryProfile: ch.camera.SecondaryProfileToken,
		},
	}
	return &a, nil
}

func (m *Manager) GetStreamUri(id model.CameraID, profile model.Profile, transport rms_cctv.VideoTransport) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ch, ok := m.channels[id]
	if !ok {
		return "", fmt.Errorf("camera %d not found", id)
	}

	streamID := ch.camera.PrimaryExternalStreamID
	if profile == model.SecondaryProfile {
		streamID = ch.camera.SecondaryExternalStreamID
	}

	uri, err := m.backend.GetStreamUri(streamID, transport)
	if err != nil {
		return "", err
	}
	return uri.String(), err
}

func (m *Manager) GetReplayUri(id model.CameraID, transport rms_cctv.VideoTransport, timestamp time.Time) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ch, ok := m.channels[id]
	if !ok {
		return "", fmt.Errorf("camera %d not found", id)
	}

	uri, err := m.backend.GetReplayUri(ch.camera.ExternalArchiveID, transport, timestamp)
	if err != nil {
		return "", err
	}
	return uri.String(), err
}

func (m *Manager) GetArchive(id model.CameraID) (accessor.Archive, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ch, ok := m.channels[id]
	if !ok {
		return nil, fmt.Errorf("camera %d not found", id)
	}
	a := archiveAccessor{
		id:       ch.camera.ExternalArchiveID,
		recorder: m.backend,
	}
	return &a, nil
}

func (m *Manager) Unregister(cam *model.Camera) error {
	if cam.ExternalArchiveID != "" {
		if err := m.backend.DeleteArchive(cam.ExternalArchiveID); err != nil {
			return fmt.Errorf("cannot delete archive: %w", err)
		}
		cam.ExternalArchiveID = ""
	}
	if cam.PrimaryExternalStreamID != "" {
		if err := m.backend.DeleteStream(cam.PrimaryExternalStreamID); err != nil {
			return fmt.Errorf("cannot delete primary stream: %w", err)
		}
		cam.PrimaryExternalStreamID = ""
	}
	if cam.SecondaryExternalStreamID != "" {
		if err := m.backend.DeleteStream(cam.SecondaryExternalStreamID); err != nil {
			return fmt.Errorf("cannot delete secondary stream: %w", err)
		}
		cam.SecondaryExternalStreamID = ""
	}

	return nil
}

func (m *Manager) removeAndReturn(id model.CameraID) *model.Camera {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch, ok := m.channels[id]
	if !ok {
		return nil
	}

	ch.l.Stop()
	delete(m.channels, id)
	return ch.camera
}

func (m *Manager) Remove(id model.CameraID) error {
	cam := m.removeAndReturn(id)
	if cam == nil {
		return errors.New("camera not found")
	}

	if err := m.Unregister(cam); err != nil {
		m.l.Logf(logger.ErrorLevel, "Unregister camera %d failed: %s", id, err)
	}

	return nil
}
