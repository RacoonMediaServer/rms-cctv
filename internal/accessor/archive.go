package accessor

type Archive interface {
	StartRecording() error
	StopRecording() error
	SetQuality(quality uint) error
}
