package tracing

import "github.com/sarchlab/akita/v4/sim"

// A MessageLog is a log of messages
type MessageLog struct {
	ID           string         `json:"id"`
	Source       string         `json:"source"`
	Destination  string         `json:"destination"`
	EnqueueTime  sim.VTimeInSec `json:"enqueue_time"`
	TransmitTime sim.VTimeInSec `json:"transmit_time"`
	ReceiveTime  sim.VTimeInSec `json:"receive_time"`
	DequeueTime  sim.VTimeInSec `json:"dequeue_time"`
	Detail       interface{}    `json:"-"`
}

// NewMessageLog creates a new MessageLog with the given ID, source, and destination.
func NewMessageLog(id, source, destination string) *MessageLog {
	return &MessageLog{
		ID:          id,
		Source:      source,
		Destination: destination,
	}
}
