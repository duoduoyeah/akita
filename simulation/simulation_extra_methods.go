// This file is outdated and should be deleted later
package simulation

import (
	"fmt"

	"github.com/sarchlab/akita/v4/sim"
)

func (s *Simulation) ComponentNames() []string {
	names := make([]string, len(s.components))
	for i, c := range s.components {
		names[i] = c.Name()
	}
	return names
}

func (s *Simulation) PortNames() []string {
	names := make([]string, len(s.ports))
	for i, p := range s.ports {
		names[i] = p.Name()
	}
	return names
}

type PortConnections struct {
	portName          string
	portComponentName string
	dstPortName       []sim.RemotePort
	incomingPorts     []sim.RemotePort
	outgoingPorts     []sim.RemotePort
	// dstPortComponentName []string
}

func (p PortConnections) Print() {
	fmt.Printf("{\n  PortName: %s,\n  PortComponentName: %s,\n",
		p.portName,
		p.portComponentName)

	// Print DstPorts with 2 items per line
	fmt.Printf("  DstPorts: [")
	for i, port := range p.dstPortName {
		if i%2 == 0 && i > 0 {
			fmt.Printf("\n")
		}
		if i%2 == 0 {
			fmt.Printf("    %s", port)
		} else {
			fmt.Printf(", %s", port)
		}
	}
	if len(p.dstPortName) > 0 {
		fmt.Printf("\n")
	}
	fmt.Printf("  ],\n")

	// Print incomingPorts with 2 items per line
	fmt.Printf("  incomingPorts: [")
	for i, port := range p.incomingPorts {
		if i%2 == 0 && i > 0 {
			fmt.Printf("\n")
		}
		if i%2 == 0 {
			fmt.Printf("    %s", port)
		} else {
			fmt.Printf(", %s", port)
		}
	}
	if len(p.incomingPorts) > 0 {
		fmt.Printf("\n")
	}
	fmt.Printf("  ],\n")

	// Print outgoingPorts with 2 items per line
	fmt.Printf("  outgoingPorts: [")
	for i, port := range p.outgoingPorts {
		if i%2 == 0 && i > 0 {
			fmt.Printf("\n")
		}
		if i%2 == 0 {
			fmt.Printf("    %s", port)
		} else {
			fmt.Printf(", %s", port)
		}
	}
	if len(p.outgoingPorts) > 0 {
		fmt.Printf("\n")
	}
	fmt.Printf("  ],\n")

	// Print DstPortComponents with 2 items per line
	// fmt.Printf("  DstPortComponents: [\n")
	// for i, component := range p.dstPortComponentName {
	// 	if i%2 == 0 && i > 0 {
	// 		fmt.Printf("\n")
	// 	}
	// 	if i%2 == 0 {
	// 		fmt.Printf("    %s", component)
	// 	} else {
	// 		fmt.Printf(", %s", component)
	// 	}
	// }
	// if len(p.dstPortComponentName) > 0 {
	// 	fmt.Printf("\n")
	// }
	// fmt.Printf("  ]\n")

	fmt.Printf("}\n")
}

// return a list of sets, each set (portName, portComponentName, dstPortName, dstPortComponentName)
func (s *Simulation) GetPortConnections() []PortConnections {
	connections := make([]PortConnections, len(s.ports))

	for i, p := range s.ports {
		if p == nil {
			continue
		}
		connections[i].portName = p.Name()
		connections[i].portComponentName = p.Component().Name()
		connections[i].dstPortName = p.GetDst()
		connections[i].incomingPorts = p.GetIncomingPorts()
		connections[i].outgoingPorts = p.GetOutgoingPorts()
	}

	return connections
}
