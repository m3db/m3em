package m3db

import (
	"github.com/m3db/m3em/node"

	"github.com/m3db/m3x/instrument"
)

// Node represents a ServiceNode pointing to an M3DB process
type Node interface {
	node.ServiceNode

	// Health returns the health for this node
	Health() (NodeHealth, error)

	// Bootstrapped returns whether the node is bootstrapped
	Bootstrapped() bool

	// TODO(prateek): add more m3db service endpoints in ServiceNode
	// - query service observable properties (nowFn, detailed_status)
	// - set nowFn offset
	// - logs
	// - metrics
}

// NodeHealth provides Health information for a M3DB node
type NodeHealth struct {
	Bootstrapped bool
	Status       string
	OK           bool
}

// Options represent the knobs to control m3db.Node behavior
type Options interface {
	// Validate validates the Options
	Validate() error

	// SetInstrumentOptions sets the instrumentation options
	SetInstrumentOptions(instrument.Options) Options

	// InstrumentOptions returns the instrumentation options
	InstrumentOptions() instrument.Options

	// SetNodeOptions sets the node options
	SetNodeOptions(node.Options) Options

	// NodeOptions returns the node options
	NodeOptions() node.Options
}
