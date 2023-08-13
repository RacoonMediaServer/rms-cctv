package reactor

type setReactionsCommand struct {
	cameraId  uint32
	reactions []Reaction
}

type dropReactionsCommand struct {
	cameraId uint32
}
