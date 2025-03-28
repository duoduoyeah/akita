
# Three Examples
This document provides three examples demonstrating core concepts and usage patterns within the Akita simulation framework. 

The examples are designed to progressively introduce complexity. 

The "Cell Split" example starts with the fundamentals of event-driven simulation, using custom events and handlers to model a simple growth process. 

The "Ping" example builds upon this by introducing components, ports, and connections, demonstrating how to model request-response interactions between distinct simulation entities using explicit event scheduling for delays. 

Finally, the "Tickingping" example refines the ping scenario by utilizing `TickingComponent` to manage internal processing delays through cycle-based logic, offering an alternative approach to time management within components.

## Cell Split
This example models exponential growth similar to cell division.

The simulation involves:

1.  **An Engine:** A `sim.NewSerialEngine` is created to manage and execute events in time order.
2.  **An Event:** A `splitEvent` represents the moment a cell divides.
3.  **A Handler:** A `handler` processes `splitEvent`s. When an event occurs:
    *   It increments a counter (`count`) representing the total number of cells.
    *   It schedules two *new* `splitEvent`s to occur at random future times (1 to 2 seconds later), simulating the division into two new cells.
4.  **Simulation Flow:**
    *   The simulation starts with an initial `count` of 1.
    *   The first `splitEvent` is scheduled at a random time.
    *   The engine runs, processing events. Each processed `splitEvent` increases the cell count and potentially schedules two more events, provided their scheduled time is before a predefined `endTime`.
5.  **Result:** The simulation stops after `endTime`, and the final cell count is printed.

This example showcases how to define custom events and handlers and use the simulation engine to model time-based processes with branching behavior.

## Ping
This example models a basic "ping" mechanism where one component sends a request and the other sends back a response after a delay, allowing the initiator to measure the round-trip time.

The simulation involves:

1.  **An Engine:** A `sim.NewSerialEngine` manages event execution.
2.  **Components (`Comp`):** Two instances (`agentA`, `agentB`) are created using a `Builder`. Each `Comp` has an `OutPort` for sending messages.
3.  **Messages:**
    *   `PingReq`: Represents the ping request message.
    *   `PingRsp`: Represents the ping response message.
4.  **Events:**
    *   `StartPingEvent`: An external event scheduled to trigger the initial ping from `agentA`.
    *   `RspPingEvent`: An internal event scheduled by the receiving component (`agentB`) to delay the response.
5.  **Connection:** A `directconnection` links the `OutPort`s of `agentA` and `agentB`, enabling message transfer.
6.  **Builder:** A `Builder` struct provides a convenient way to configure and create `Comp` instances, setting the `Engine`.
7.  **Simulation Flow:**
    *   The simulation is set up with the engine, two `Comp` instances (`agentA`, `agentB`), and a connection between them.
    *   `StartPingEvent`s are scheduled for `agentA` at specific times (e.g., 1s and 3s).
    *   When `agentA` handles a `StartPingEvent`:
        *   It creates and sends a `PingReq` message to `agentB` via its `OutPort` and the connection.
        *   It records the start time associated with the request's sequence ID.
    *   The connection delivers the `PingReq` to `agentB`. `agentB`'s `NotifyRecv` method is called.
    *   `agentB` processes the `PingReq` by scheduling an `RspPingEvent` to occur 2 seconds later.
    *   When `agentB` handles the `RspPingEvent`:
        *   It creates and sends a `PingRsp` message back to `agentA`.
    *   The connection delivers the `PingRsp` to `agentA`. `agentA`'s `NotifyRecv` method is called.
    *   `agentA` processes the `PingRsp` by retrieving the original start time (using the sequence ID) and calculating the total round-trip duration.
    *   The sequence ID and the calculated duration are printed.
8.  **Result:** The simulation prints the round-trip time for each ping initiated.

This example showcases:
*   Defining custom message types (`PingReq`, `PingRsp`).
*   Component interaction via ports and connections.
*   Handling incoming messages using `NotifyRecv`.
*   Using internal events (`RspPingEvent`) to model processing delays.
*   Using external events (`StartPingEvent`) to initiate actions.
*   Measuring time within the simulation.
*   Employing a Builder pattern for component setup.

### New Concepts in Ping

* **Components:** Instead of a single handler, the simulation uses distinct `Component` instances (`agentA`, `agentB`) representing separate entities.
* **Ports:** Components communicate through defined interfaces called `Port`s (`OutPort`).
* **Connections:** `Connection`s (`directconnection`) are used to link ports and facilitate message transfer between components.
* **Messages:** Data exchanged between components is structured as `Msg` types (`PingReq`, `PingRsp`), rather than just triggering events.
* **Component Interaction Model:** It demonstrates how components send messages through ports and receive them via the `NotifyRecv` mechanism, triggered by the connection delivering a message.
* **Builder Pattern:** Introduces a structured way to configure and instantiate components.

## Tickingping
This example revisits the ping-pong scenario but utilizes Akita's `TickingComponent` instead of explicit event scheduling for internal processing delays. It demonstrates how component logic can be driven by discrete time steps (ticks) based on a defined frequency.

The simulation involves:

1.  **An Engine:** A `sim.NewSerialEngine` manages the simulation time and event queue.
2.  **Components (`Comp`):**
    *   Two instances (`agentA`, `agentB`) are created using a `Builder`.
    *   Each `Comp` embeds `sim.TickingComponent`, meaning its primary logic executes within a `Tick()` function called periodically by the engine based on the component's frequency.
    *   It uses `sim.MiddlewareHolder` and a custom `middleware` struct to structure the logic performed during each tick.
    *   Each `Comp` has an `OutPort` for sending messages.
3.  **Messages:**
    *   `PingReq`: Represents the ping request message (same as in the Ping example).
    *   `PingRsp`: Represents the ping response message (same as in the Ping example).
4.  **Middleware:** The `middleware` struct contains the core logic:
    *   `processInput()`: Checks the `OutPort` for incoming messages (`PingReq` or `PingRsp`) at the start of a tick.
    *   `processingPingReq()`: When a `PingReq` is received, it initiates a `pingTransaction` which includes a `cycleLeft` counter (e.g., 2 cycles) to simulate processing time.
    *   `processingPingRsp()`: When a `PingRsp` is received, it calculates the round-trip time and prints it.
    *   `countDown()`: Decrements the `cycleLeft` counter for active transactions each tick.
    *   `sendRsp()`: Sends a `PingRsp` when a transaction's `cycleLeft` reaches zero.
    *   `sendPing()`: Sends a new `PingReq` if required (`numPingNeedToSend > 0`).
5.  **Connection:** A `directconnection` links the `OutPort`s of `agentA` and `agentB`.
6.  **Builder:** A `Builder` struct configures and creates `Comp` instances, setting the `engine` and operating `freq` (frequency).
7.  **Simulation Flow:**
    *   The simulation is set up with the engine, two `Comp` instances (`agentA`, `agentB` operating at 1 Hz), and a connection.
    *   `agentA` is configured to send 2 pings (`numPingNeedToSend = 2`) to `agentB`.
    *   `agentA.TickLater()` schedules the first tick for `agentA`.
    *   The engine runs.
    *   On its scheduled ticks, `agentA`'s `middleware.Tick()`:
        *   Calls `sendPing()` to send `PingReq` messages (one per tick until `numPingNeedToSend` is 0), recording the start time.
        *   Calls `processInput()` to check for incoming `PingRsp`. If found, `processingPingRsp()` calculates and prints the duration.
    *   The connection delivers `PingReq`s to `agentB`.
    *   On its scheduled ticks, `agentB`'s `middleware.Tick()`:
        *   Calls `processInput()` to check for incoming `PingReq`. If found, `processingPingReq()` starts a transaction with a 2-cycle delay.
        *   Calls `countDown()` to decrement the delay counter for active transactions.
        *   Calls `sendRsp()` to send a `PingRsp` back to `agentA` once a transaction's delay counter reaches zero.
    *   The connection delivers `PingRsp`s back to `agentA`.
    *   The simulation continues until no component can make further progress in a tick.
8.  **Result:** The simulation prints the round-trip time for each ping, reflecting the processing delay (2 cycles = 2 seconds at 1 Hz) within `agentB` plus any potential queuing or transmission delays.

This example showcases:
*   Using `TickingComponent` for cycle-driven behavior.
*   Implementing component logic within a `Tick()` function, often structured using middleware.
*   Modeling processing delays using cycle counters within the `Tick()` logic.
*   Interaction between ticking components via ports and connections.

### New Concepts in Tickingping

Compared to the basic `Ping` example, `Tickingping` introduces several key concepts:

*   **Execution Model:** `Ping` relies on explicitly scheduled `sim.Event`s to trigger actions (sending a ping, responding). `Tickingping` uses the `sim.TickingComponent` model, where component logic executes periodically within a `Tick()` function based on the component's frequency. This aligns better with hardware components operating on clock cycles.
*   **Internal State Management:** `Tickingping` manages internal state (like processing delays or pending requests) within the component's `Tick()` function and associated data structures (like the `middleware` and `pingTransaction`). `Ping` handles state transitions primarily through event handlers.
*   **Modeling Processing Time:** `Tickingping` explicitly models the time taken for internal processing using a cycle counter (`cycleLeft`) within the `Tick()` function. The basic `Ping` example doesn't model this internal component delay; delays primarily arise from message transmission through the connection.
*   **Component Structure:** `Tickingping` utilizes `sim.TickingComponent` and often employs patterns like middleware (`sim.MiddlewareHolder`) to organize the logic executed during each tick. `Ping` uses a simpler `sim.ComponentBase` structure.
*   **Initiation:** While `Ping` uses specific events (`StartPingEvent`) scheduled externally to start the process, `Tickingping` often initiates actions from within its `Tick()` function based on its internal state (e.g., `numPingNeedToSend`) after an initial `TickLater()` call.
