package event_bus

type BaseEvent struct {
	AccountID string
}

type PeerCreated struct {
	BaseEvent
	PeerID string
}

type PortAllocated struct {
	BaseEvent
	PeerID string
	Port   int
}
