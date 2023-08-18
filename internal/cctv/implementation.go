package cctv

import (
	"net/url"
)

type backend struct {
}

func (b backend) AddStream(streamUrl *url.URL) (ID, error) {
	//TODO implement me
	panic("implement me")
}

func (b backend) DeleteStream(id ID) error {
	//TODO implement me
	panic("implement me")
}

func (b backend) GetStreamUri(id ID) (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (b backend) StartRecording(id ID) error {
	//TODO implement me
	panic("implement me")
}

func (b backend) StopRecording(id ID) error {
	//TODO implement me
	panic("implement me")
}

func (b backend) SetQuality(id ID, quality uint) error {
	//TODO implement me
	panic("implement me")
}

func (b backend) AddArchive(id ID, rotationDays uint) (ID, error) {
	//TODO implement me
	panic("implement me")
}

func (b backend) DeleteArchive(id ID) error {
	//TODO implement me
	panic("implement me")
}
