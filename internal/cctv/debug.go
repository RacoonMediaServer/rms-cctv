package cctv

import (
	"errors"
	"fmt"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"go-micro.dev/v4/logger"
	"net/url"
	"sync"
	"time"
)

type channel struct {
	u *url.URL
}

type archive struct {
	u         *url.URL
	ch        ID
	recording bool
	quality   uint
}

type debugBackend struct {
	mu             sync.Mutex
	l              logger.Logger
	channels       map[ID]*channel
	archives       map[ID]*archive
	streamCounter  int
	archiveCounter int
}

func (b *debugBackend) AddStream(streamUrl *url.URL) (ID, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := channel{u: streamUrl}
	b.streamCounter++
	id := ID(fmt.Sprintf("ch%d", b.streamCounter))

	b.channels[id] = &ch
	b.l.Logf(logger.InfoLevel, "Stream '%s' [ %s ] registered", id, streamUrl)

	return id, nil
}

func (b *debugBackend) DeleteStream(id ID) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.channels[id]
	if !ok {
		return errors.New("not found")
	}

	delete(b.channels, id)
	b.l.Logf(logger.InfoLevel, "Stream '%s' deleted", id)

	return nil
}

func (b *debugBackend) GetStreamUri(id ID, transport rms_cctv.VideoTransport) (*url.URL, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch, ok := b.channels[id]
	if !ok {
		return nil, errors.New("not found")
	}

	return ch.u, nil
}

func (b *debugBackend) StartRecording(id ID) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	a, ok := b.archives[id]
	if !ok {
		return errors.New("not found")
	}

	b.l.Logf(logger.InfoLevel, "[%s of %s] Recording %t => %t", id, a.ch, a.recording, true)
	a.recording = true

	return nil
}

func (b *debugBackend) StopRecording(id ID) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	a, ok := b.archives[id]
	if !ok {
		return errors.New("not found")
	}

	b.l.Logf(logger.InfoLevel, "[%s of %s] Recording %t => %t", id, a.ch, a.recording, false)
	a.recording = false

	return nil
}

func (b *debugBackend) SetQuality(id ID, quality uint) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	a, ok := b.archives[id]
	if !ok {
		return errors.New("not found")
	}

	b.l.Logf(logger.InfoLevel, "[%s of %s] Quality %d => %d", id, a.ch, a.recording, quality)
	a.quality = quality

	return nil
}

func (b *debugBackend) AddArchive(streamID ID, rotationDays uint) (ID, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch, ok := b.channels[streamID]
	if !ok {
		return "", errors.New("stream not found")
	}

	b.archiveCounter++
	id := ID(fmt.Sprintf("rec%d", b.archiveCounter))

	b.archives[id] = &archive{u: ch.u, ch: streamID}
	b.l.Logf(logger.InfoLevel, "Archive '%s' [ %s ] registered", id, streamID)

	return id, nil
}

func (b *debugBackend) DeleteArchive(id ID) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.archives[id]
	if !ok {
		return errors.New("not found")
	}

	delete(b.archives, id)
	b.l.Logf(logger.InfoLevel, "Archive '%s' deleted", id)

	return nil
}

func (b *debugBackend) GetReplayUri(id ID, transport rms_cctv.VideoTransport, timestamp time.Time) (*url.URL, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	a, ok := b.archives[id]
	if !ok {
		return nil, errors.New("not found")
	}

	return a.u, nil
}
