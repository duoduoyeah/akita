package ping

import (
	"github.com/sarchlab/akita/v4/sim"
	"github.com/sarchlab/akita/v4/sim/directconnection"
)

// Example_pingWithEvents demonstrates how to use the ping component with events.
//
// 1. It creates a simulation engine.
// 2. It creates two ping agents, agentA and agentB.
// 3. It creates a direct connection between the two agents.
// 4. It plugs the agents' ports into the connection.
// 5. It creates two StartPingEvent events to initiate pings from agentA to agentB.
// 6. It schedules the events.
// 7. It runs the simulation.
func Example_pingWithEvents() {
	engine := sim.NewSerialEngine()
	// agentA := NewPingAgent("AgentA", engine)
	agentA := MakeBuilder().WithEngine(engine).Build("AgentA")
	// agentB := NewPingAgent("AgentB", engine)
	agentB := MakeBuilder().WithEngine(engine).Build("AgentB")
	conn := directconnection.MakeBuilder().
		WithEngine(engine).
		WithFreq(1 * sim.GHz).
		Build("Conn")

	conn.PlugIn(agentA.OutPort)
	conn.PlugIn(agentB.OutPort)

	e1 := StartPingEvent{
		EventBase: sim.NewEventBase(1, agentA),
		Dst:       agentB.OutPort.AsRemote(),
	}
	e2 := StartPingEvent{
		EventBase: sim.NewEventBase(3, agentA),
		Dst:       agentB.OutPort.AsRemote(),
	}
	engine.Schedule(e1)
	engine.Schedule(e2)

	engine.Run()
	// Output:
	// Ping 0, 2.00
	// Ping 1, 2.00
}
