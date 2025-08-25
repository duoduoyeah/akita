package tracing

import (
	"sync"

	"github.com/sarchlab/akita/v4/datarecording"
	"github.com/sarchlab/akita/v4/sim"
	"github.com/tebeka/atexit"
)

type messageTableEntry struct {
	ID           string  `json:"id" akita_data:"unique"`
	Source       string  `json:"source" akita_data:"index"`
	Destination  string  `json:"destination" akita_data:"index"`
	EnqueueTime  float64 `json:"enqueue_time" akita_data:"index"`
	TransmitTime float64 `json:"transmit_time" akita_data:"index"`
	ReceiveTime  float64 `json:"receive_time" akita_data:"index"`
	DequeueTime  float64 `json:"dequeue_time" akita_data:"index"`
}

type TopologyPortEntry struct {
	Port      string `json:"port" akita_data:"unique"`
	Component string `json:"component" akita_data:"index"`
}

type PortConnectionEntry struct {
	SourcePort      string `json:"from_port" akita_data:"index"` // Source port ID
	DestinationPort string `json:"to_port" akita_data:"index"`   // Destination port ID
}

type MessageTracer struct {
	mu         sync.Mutex
	timeTeller sim.TimeTeller
	backend    datarecording.DataRecorder

	tracingMessages map[string]MessageLog
}

func NewMessageTracer(
	timeTeller sim.TimeTeller,
	dataRecorder datarecording.DataRecorder,
) *MessageTracer {
	dataRecorder.CreateTable("message_trace", messageTableEntry{})
	dataRecorder.CreateTable("topology_ports", TopologyPortEntry{})
	dataRecorder.CreateTable("ports_connection", PortConnectionEntry{})

	t := &MessageTracer{
		timeTeller:      timeTeller,
		backend:         dataRecorder,
		tracingMessages: make(map[string]MessageLog),
	}

	atexit.Register(func() {
		t.Terminate()
	})

	return t
}

// EnqueueMessage marks the enqueue of a message.
func (t *MessageTracer) EnqueueMessage(rawMsg sim.Msg) {
	t.mu.Lock()
	defer t.mu.Unlock()

	msg := NewMessageLog(rawMsg.Meta().ID, string(rawMsg.Meta().Src), string(rawMsg.Meta().Dst))

	msg.EnqueueTime = t.timeTeller.CurrentTime()

	existingMsg, found := t.tracingMessages[msg.ID]
	if !found {
		t.tracingMessages[msg.ID] = *msg
		return
	}

	existingMsg.Source = msg.Source
	existingMsg.Destination = msg.Destination
	existingMsg.EnqueueTime = msg.EnqueueTime

	t.tracingMessages[msg.ID] = existingMsg
}

// TransmitMessage marks the transmit time of a message.
func (t *MessageTracer) TransmitMessage(rawMsg sim.Msg) {
	t.mu.Lock()
	defer t.mu.Unlock()

	msgID := rawMsg.Meta().ID
	msg, found := t.tracingMessages[msgID]
	if !found {
		return
	}

	msg.TransmitTime = t.timeTeller.CurrentTime()
	t.tracingMessages[msgID] = msg
}

// ReceiveMessage marks the receive time of a message.
func (t *MessageTracer) ReceiveMessage(rawMsg sim.Msg) {
	t.mu.Lock()
	defer t.mu.Unlock()

	msgID := rawMsg.Meta().ID
	msg, found := t.tracingMessages[msgID]
	if !found {
		return
	}

	msg.ReceiveTime = t.timeTeller.CurrentTime()
	t.tracingMessages[msgID] = msg
}

// DequeueMessage marks the dequeue time of a message.
func (t *MessageTracer) DequeueMessage(rawMsg sim.Msg) {
	t.mu.Lock()
	defer t.mu.Unlock()

	msgID := rawMsg.Meta().ID
	msg, found := t.tracingMessages[msgID]
	if !found {
		return
	}

	msg.DequeueTime = t.timeTeller.CurrentTime()

	entry := messageTableEntry{
		ID:           msg.ID,
		Source:       msg.Source,
		Destination:  msg.Destination,
		EnqueueTime:  float64(msg.EnqueueTime),
		TransmitTime: float64(msg.TransmitTime),
		ReceiveTime:  float64(msg.ReceiveTime),
		DequeueTime:  float64(msg.DequeueTime),
	}
	t.backend.InsertData("message_trace", entry)

	delete(t.tracingMessages, msgID)
}

func (t *MessageTracer) Terminate() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, msg := range t.tracingMessages {
		entry := messageTableEntry{
			ID:           msg.ID,
			Source:       msg.Source,
			Destination:  msg.Destination,
			EnqueueTime:  float64(msg.EnqueueTime),
			TransmitTime: float64(msg.TransmitTime),
			ReceiveTime:  float64(msg.ReceiveTime),
			DequeueTime:  float64(msg.DequeueTime),
		}
		t.backend.InsertData("message_trace", entry)
	}

	t.tracingMessages = nil
	t.backend.Flush()
}

// Add topology port map to the topology_ports table
func (t *MessageTracer) AddTopologyPortMap(components []sim.Component) {
	for _, component := range components {
		componentName := component.Name()
		ports := component.Ports()

		for _, port := range ports {
			// Insert port info into topology_ports table
			portEntry := TopologyPortEntry{
				Port:      port.Name(),
				Component: componentName,
			}
			t.backend.InsertData("topology_ports", portEntry)

			// Insert port connections into ports_connection table
			for _, incoming := range port.GetIncomingPorts() {
				connEntry := PortConnectionEntry{
					SourcePort:      string(incoming),
					DestinationPort: port.Name(),
				}
				t.backend.InsertData("ports_connection", connEntry)
			}
			for _, outgoing := range port.GetOutgoingPorts() {
				connEntry := PortConnectionEntry{
					SourcePort:      port.Name(),
					DestinationPort: string(outgoing),
				}
				t.backend.InsertData("ports_connection", connEntry)
			}
		}
	}
}
