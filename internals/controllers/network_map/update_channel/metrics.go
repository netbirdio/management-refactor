package update_channel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type metrics struct {
	createChannelDurationMicro        metric.Int64Histogram
	closeChannelDurationMicro         metric.Int64Histogram
	closeChannelsDurationMicro        metric.Int64Histogram
	closeChannels                     metric.Int64Histogram
	sendUpdateDurationMicro           metric.Int64Histogram
	getAllConnectedPeersDurationMicro metric.Int64Histogram
	getAllConnectedPeers              metric.Int64Histogram
	hasChannelDurationMicro           metric.Int64Histogram
}

func newMetrics(meter metric.Meter) (*metrics, error) {
	createChannelDurationMicro, err := meter.Int64Histogram("management.updatechannel.create.duration.micro",
		metric.WithUnit("microseconds"),
		metric.WithDescription("Duration of how long it takes to create a new peer update channel"),
	)
	if err != nil {
		return nil, err
	}

	closeChannelDurationMicro, err := meter.Int64Histogram("management.updatechannel.close.one.duration.micro",
		metric.WithUnit("microseconds"),
		metric.WithDescription("Duration of how long it takes to close a peer update channel"),
	)
	if err != nil {
		return nil, err
	}

	closeChannelsDurationMicro, err := meter.Int64Histogram("management.updatechannel.close.multiple.duration.micro",
		metric.WithUnit("microseconds"),
		metric.WithDescription("Duration of how long it takes to close a set of peer update channels"),
	)

	if err != nil {
		return nil, err
	}

	closeChannels, err := meter.Int64Histogram("management.updatechannel.close.multiple.channels",
		metric.WithUnit("1"),
		metric.WithDescription("Number of peer update channels that have been closed"),
	)

	if err != nil {
		return nil, err
	}

	sendUpdateDurationMicro, err := meter.Int64Histogram("management.updatechannel.send.duration.micro",
		metric.WithUnit("microseconds"),
		metric.WithDescription("Duration of how long it takes to send an network map update to a peer"),
	)
	if err != nil {
		return nil, err
	}

	getAllConnectedPeersDurationMicro, err := meter.Int64Histogram("management.updatechannel.get.all.duration.micro",
		metric.WithUnit("microseconds"),
		metric.WithDescription("Duration of how long it takes to get all connected peers"),
	)
	if err != nil {
		return nil, err
	}

	getAllConnectedPeers, err := meter.Int64Histogram("management.updatechannel.get.all.peers",
		metric.WithUnit("1"),
		metric.WithDescription("Number of connected peers"),
	)
	if err != nil {
		return nil, err
	}

	hasChannelDurationMicro, err := meter.Int64Histogram("management.updatechannel.haschannel.duration.micro",
		metric.WithUnit("microseconds"),
		metric.WithDescription("Duration of how long it takes to check if a peer has a channel"),
	)
	if err != nil {
		return nil, err
	}

	return &metrics{
		createChannelDurationMicro:        createChannelDurationMicro,
		closeChannelDurationMicro:         closeChannelDurationMicro,
		closeChannelsDurationMicro:        closeChannelsDurationMicro,
		closeChannels:                     closeChannels,
		sendUpdateDurationMicro:           sendUpdateDurationMicro,
		getAllConnectedPeersDurationMicro: getAllConnectedPeersDurationMicro,
		getAllConnectedPeers:              getAllConnectedPeers,
		hasChannelDurationMicro:           hasChannelDurationMicro,
	}, nil
}

// CountCreateChannelDuration counts the duration of the CreateChannel method,
// closed indicates if existing channel was closed before creation of a new one
func (m *metrics) CountCreateChannelDuration(duration time.Duration, closed bool) {
	opts := metric.WithAttributeSet(attribute.NewSet(attribute.Bool("closed", closed)))
	m.createChannelDurationMicro.Record(context.Background(), duration.Microseconds(), opts)
}

// CountCloseChannelDuration counts the duration of the CloseChannel method
func (m *metrics) CountCloseChannelDuration(duration time.Duration) {
	m.closeChannelDurationMicro.Record(context.Background(), duration.Microseconds())
}

// CountCloseChannelsDuration counts the duration of the CloseChannels method and the number of channels have been closed
func (m *metrics) CountCloseChannelsDuration(duration time.Duration, channels int) {
	m.closeChannelsDurationMicro.Record(context.Background(), duration.Microseconds())
	m.closeChannels.Record(context.Background(), int64(channels))
}

// CountSendUpdateDuration counts the duration of the SendUpdate method
// found indicates if peer had channel, dropped indicates if the message was dropped due channel buffer overload
func (m *metrics) CountSendUpdateDuration(duration time.Duration, found, dropped bool) {
	opts := metric.WithAttributeSet(attribute.NewSet(attribute.Bool("found", found), attribute.Bool("dropped", dropped)))
	m.sendUpdateDurationMicro.Record(context.Background(), duration.Microseconds(), opts)
}

// CountGetAllConnectedPeersDuration counts the duration of the GetAllConnectedPeers method and the number of peers have been returned
func (m *metrics) CountGetAllConnectedPeersDuration(duration time.Duration, peers int) {
	m.getAllConnectedPeersDurationMicro.Record(context.Background(), duration.Microseconds())
	m.getAllConnectedPeers.Record(context.Background(), int64(peers))
}

// CountHasChannelDuration counts the duration of the HasChannel method
func (m *metrics) CountHasChannelDuration(duration time.Duration) {
	m.hasChannelDurationMicro.Record(context.Background(), duration.Microseconds())
}
