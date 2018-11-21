package gm

type OperationType int

const (
	Update OperationType = iota + 1
)

type Operation struct {
	repo Repository
	opType OperationType
}

func (op Operation) Execute() {
	switch op.opType {
	case Update:
		op.repo.Update()
	}
}

func (gmc *GitManagerConfig) Loop() {
	gmc.ch = make(chan Operation, 2)

	for {
		select {
			case op, ok := <-gmc.ch:
			if !ok {
				return
			}
			op.Execute()
		}
	}
}
