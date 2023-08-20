package cctv

import (
	"net/url"
)

type backend struct {
}

func (b backend) AddStream(streamUrl *url.URL) (ID, error) {
	//TODO implement me
	return "ch1", nil
}

func (b backend) DeleteStream(id ID) error {
	//TODO implement me
	return nil
}

func (b backend) GetStreamUri(id ID) (*url.URL, error) {
	//TODO implement me
	return url.Parse("rtsp://127.0.0.1/stream1")
}

func (b backend) StartRecording(id ID) error {
	//TODO implement me
	return nil
}

func (b backend) StopRecording(id ID) error {
	//TODO implement me
	return nil
}

func (b backend) SetQuality(id ID, quality uint) error {
	//TODO implement me
	return nil
}

func (b backend) AddArchive(id ID, rotationDays uint) (ID, error) {
	//TODO implement me
	return "rec1", nil
}

func (b backend) DeleteArchive(id ID) error {
	//TODO implement me
	return nil
}
