package operation

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

type Operation struct {
	Type   OperationType
	Local  string
	Remote string
}

type Operations []Operation
