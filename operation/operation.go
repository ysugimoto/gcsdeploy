package operation

import (
	"github.com/ysugimoto/gcsdeploy/local"
	"github.com/ysugimoto/gcsdeploy/remote"
)

type OperationType int

const (
	Add OperationType = iota + 1
	Update
	Delete
)

func (o OperationType) String() string {
	switch o {
	case Add:
		return "ADD"
	case Update:
		return "UPDATE"
	case Delete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

// Operation struct holds all informations to do the GCS.
// Operation type, local object, and remote object.
type Operation struct {
	Type   OperationType
	Local  local.Object
	Remote remote.Object
}

type Operations []Operation
