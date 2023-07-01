package transaction

import "fmt"

type TrStatus int

const (
	Active TrStatus = iota
	Finished
	Aborted
	Waiting
)

func (tr TrStatus) String() string {
	switch tr {
	case Active:
		return "Ativa"
	case Finished:
		return "Concluída"
	case Aborted:
		return "Abortada"
	default:
		return "Esperando"
	}
}

type TrManager struct {
	transactions map[int]*Transaction
}

func NewTrManager() *TrManager {
	return &TrManager{make(map[int]*Transaction)}
}

func NewTrManagerFromOperationsTable(opTable OperationsTable) *TrManager {
	trManager := NewTrManager()

	for id, operations := range opTable {
		trManager.AddTransaction(NewTransaction(id, operations))
	}

	return trManager
}

func (trm *TrManager) AddTransaction(tr *Transaction) {
	if _, ok := trm.transactions[tr.ID()]; !ok {
		trm.transactions[tr.ID()] = tr
		return
	}
}

func (trm *TrManager) Transactions() map[int]*Transaction {
	return trm.transactions
}

func (trm *TrManager) GetTransactionStatus(id int) (TrStatus, error) {
	if tr, ok := trm.transactions[id]; ok {
		return tr.Status(), nil
	}

	return Finished, nil
}

func (trm *TrManager) UpdateStatus(id int, status TrStatus) {
	if tr, ok := trm.transactions[id]; ok {
		tr.UpdateStatus(status)
	}
}

func (trm *TrManager) PrintTransactions() {
	fmt.Println("|------------------------|")
	fmt.Println("|   TrID   |   Status    |")
	fmt.Println("|------------------------|")
	for id, transaction := range trm.transactions {
		fmt.Printf("|   %-6d |   %-9s |\n", id, transaction.Status())
	}
	fmt.Println("|------------------------|")
}
