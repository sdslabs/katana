package deployment

type TeamConfig struct {
	TeamCount uint
	TeamLabel string
}

type BroadcastConfig struct {
	BroadcastCount uint
	BroadcastLabel string
	BroadcastPort  uint
}

type BroadcastServiceConfig struct {
	BroadcastLabel string
	BroadcastPort  uint
}
