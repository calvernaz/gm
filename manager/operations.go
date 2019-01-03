package manager

import (
	"log"
)

type OperationType int

// operation type
const (
	Update OperationType = iota + 1
	Download
)

type Operation struct {
	Repo   RepositoryEntry
	OpType OperationType
}

func (op Operation) Execute() error {
	switch op.OpType {
	case Update:
		err := op.Repo.Update()
		return err
	case Download:
		return op.Repo.Download()
	default:
		log.Println("no valid operation")
		return nil
	}
}
