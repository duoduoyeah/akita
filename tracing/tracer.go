package tracing

import "github.com/sarchlab/akita/v4/sim"

// A Tracer can collect task traces
type Tracer interface {
	StartTask(task Task)
	StepTask(task Task)
	AddMilestone(milestone Milestone)
	EndTask(task Task)
}

type MsgTracer interface {
	EnqueueMessage(msg sim.Msg)
	TransmitMessage(msg sim.Msg)
	ReceiveMessage(msg sim.Msg)
	DequeueMessage(msg sim.Msg)
}
