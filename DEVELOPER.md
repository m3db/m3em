Developer Notes
===============

Utility Packages:
  `tcp` - utility package *adopted* from go-common (private; tiny)
  `build` - utility package to represent build/confs (public; tiny)
  `os` - utility package for FileIteration and Process Lifecycle management (private)
  `cluster` - is a collection of instances along with corresponding placement alteration (public)

Packages for client<->server communications, life-cycle maintenance, event propagation.
  `environment` - main construct is ServiceNode, which wraps operator.Operator and services.PlacementInstance (public)
  `operator` - the client side Operator implementation, server side Heartbeat implementation (public)
  `generated` - grpc schemas for Operator, and Heartbeat services (private)
  `agent` - the server side Operator implementation, client side Heartbeat implementation (private)
  `services` - m3em_agent main (private)
  `integration` - currently has tests for file transfer, process lifecycle, and heartbeating; I plan to add cluster tests as well (private)

Things to be aware of:
- I'm planning on merging operator.Operator and ServiceNode to reduce the exposed API. Any objections?
- ServiceNode only is specific to M3 only because of the Health check it exposes. So I'm going to remove those bits into another package and rename ServiceNode to 'ServiceNode' or some such
- Currently the os processManager interface does a "fork/exec" using go-routines, would it be worthwhile to investigate having the process run in an external process manager and use APIs to interact with it instead. This would buy us the ability to isolate the agent process lifecycle from the test process lifecycle. It might be future work if we want to use m3em agents as long running processes used to deploy m3db to do this but wanted to hear people's thoughts on it.