package transaction

import (
	"fmt"
)

type Transaction struct {
	id         int
	status     TrStatus
	operations []*Operation
}

func NewTransaction(id int, operations []*Operation) *Transaction {
	return &Transaction{id, Active, operations}
}

func (tm *Transaction) ID() int {
	return tm.id
}

func (tm *Transaction) Status() TrStatus {
	return tm.status
}

func (tm *Transaction) Operation() []*Operation {
	return tm.operations
}

func (tm *Transaction) UpdateStatus(status TrStatus) {
	tm.status = status
}

func (tm *Transaction) String() string {
	str := "[ "

	for _, op := range tm.operations {
		str += op.String() + " "
	}
	str += "]"

	return fmt.Sprintf("(ID: %d, STATUS: %s, OPERATIONS: %s)", tm.id, tm.status, str)
}
