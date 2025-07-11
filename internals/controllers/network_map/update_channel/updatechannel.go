package update_channel

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/management-refactor/internals/controllers/network_map"
	appmetrics "github.com/netbirdio/management-refactor/internals/shared/metrics"
)

const channelBufferSize = 100

type UpdateChannel struct {
	// peerChannels is an update channel indexed by Peer.ID
	peerChannels map[string]chan *network_map.UpdateMessage
	// channelsMux keeps the mutex to access peerChannels
	channelsMux *sync.RWMutex
	// metrics provides method to collect application metrics
	metrics *metrics
}

// NewUpdateChannel returns a new instance of UpdateChannel
func NewUpdateChannel(metrics *appmetrics.AppMetrics) *UpdateChannel {
	cMetrics, err := appmetrics.RegisterMetrics(metrics, newMetrics)
	if err != nil {
		log.Fatalf("Failed to register updatechannel metrics: %v", err)
	}
	return &UpdateChannel{
		peerChannels: make(map[string]chan *network_map.UpdateMessage),
		channelsMux:  &sync.RWMutex{},
		metrics:      cMetrics,
	}
}

// SendUpdate sends update message to the peer's channel
func (p *UpdateChannel) SendUpdate(ctx context.Context, peerID string, update *network_map.UpdateMessage) {
	log.Debugf("Sending update message to peer %s on public update channel", peerID)

	start := time.Now()
	var found, dropped bool

	p.channelsMux.RLock()

	defer func() {
		p.channelsMux.RUnlock()
		p.metrics.CountSendUpdateDuration(time.Since(start), found, dropped)
	}()

	if channel, ok := p.peerChannels[peerID]; ok {
		found = true
		select {
		case channel <- update:
			log.WithContext(ctx).Debugf("update was sent to channel for peer %s", peerID)
		default:
			dropped = true
			log.WithContext(ctx).Warnf("channel for peer %s is %d full or closed", peerID, len(channel))
		}
	} else {
		log.WithContext(ctx).Debugf("peer %s has no channel", peerID)
	}
}

// CreateChannel creates a go channel for a given peer used to deliver updates relevant to the peer.
func (p *UpdateChannel) CreateChannel(ctx context.Context, peerID string) chan *network_map.UpdateMessage {
	start := time.Now()

	closed := false

	p.channelsMux.Lock()
	defer func() {
		p.channelsMux.Unlock()
		p.metrics.CountCreateChannelDuration(time.Since(start), closed)
	}()

	if channel, ok := p.peerChannels[peerID]; ok {
		closed = true
		delete(p.peerChannels, peerID)
		close(channel)
	}
	// mbragin: todo shouldn't it be more? or configurable?
	channel := make(chan *network_map.UpdateMessage, channelBufferSize)
	p.peerChannels[peerID] = channel

	log.WithContext(ctx).Debugf("opened updates channel for a peer %s", peerID)

	return channel
}

func (p *UpdateChannel) closeChannel(ctx context.Context, peerID string) {
	if channel, ok := p.peerChannels[peerID]; ok {
		delete(p.peerChannels, peerID)
		close(channel)

		log.WithContext(ctx).Debugf("closed updates channel of a peer %s", peerID)
		return
	}

	log.WithContext(ctx).Debugf("closing updates channel: peer %s has no channel", peerID)
}

// CloseChannels closes updates channel for each given peer
func (p *UpdateChannel) CloseChannels(ctx context.Context, peerIDs []string) {
	start := time.Now()

	p.channelsMux.Lock()
	defer func() {
		p.channelsMux.Unlock()
		p.metrics.CountCloseChannelsDuration(time.Since(start), len(peerIDs))
	}()

	for _, id := range peerIDs {
		p.closeChannel(ctx, id)
	}
}

// CloseChannel closes updates channel of a given peer
func (p *UpdateChannel) CloseChannel(ctx context.Context, peerID string) {
	start := time.Now()

	p.channelsMux.Lock()
	defer func() {
		p.channelsMux.Unlock()
		p.metrics.CountCloseChannelDuration(time.Since(start))
	}()

	p.closeChannel(ctx, peerID)
}

// GetAllConnectedPeers returns a copy of the connected peers map
func (p *UpdateChannel) GetAllConnectedPeers() map[string]struct{} {
	start := time.Now()

	p.channelsMux.RLock()

	m := make(map[string]struct{})

	defer func() {
		p.channelsMux.RUnlock()
		p.metrics.CountGetAllConnectedPeersDuration(time.Since(start), len(m))
	}()

	for ID := range p.peerChannels {
		m[ID] = struct{}{}
	}

	return m
}

// HasChannel returns true if peers has channel in update manager, otherwise false
func (p *UpdateChannel) HasChannel(peerID string) bool {
	start := time.Now()

	p.channelsMux.RLock()

	defer func() {
		p.channelsMux.RUnlock()
		// TODO; remove condition as it is useless
		if p.metrics != nil {
			p.metrics.CountHasChannelDuration(time.Since(start))
		}
	}()

	_, ok := p.peerChannels[peerID]

	return ok
}
