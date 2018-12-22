package gm

type OperationType int

// operation type
const (
	Update OperationType = iota + 1
)

const (
	Success bool = true
)

type Operation struct {
	Repo Repository
	OpType OperationType
}

func (op Operation) Execute() {
	switch op.OpType {
	case Update:
		if err := op.Repo.Update(); err != nil {
		
		}
	}
}

