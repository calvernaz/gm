package gm

type OperationType int

const (
	Update OperationType = iota + 1
)

type Operation struct {
	Repo Repository
	OpType OperationType
}

func (op Operation) Execute() error {
	switch op.OpType {
	case Update:
		return op.Repo.Update()
	}
	return nil
}

