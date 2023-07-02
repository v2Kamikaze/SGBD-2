package transaction

import (
	"fmt"
)

type Transaction struct {
	id     int
	status TrStatus
}

func NewTransaction(id int) *Transaction {
	return &Transaction{id, Active}
}

func (tm *Transaction) ID() int {
	return tm.id
}

func (tm *Transaction) Status() TrStatus {
	return tm.status
}

func (tm *Transaction) UpdateStatus(status TrStatus) {
	tm.status = status
}

func (tm *Transaction) String() string {
	return fmt.Sprintf("(ID: %d, STATUS: %s)", tm.id, tm.status)
}
