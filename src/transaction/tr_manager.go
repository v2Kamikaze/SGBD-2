package transaction

import (
	"fmt"
)

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
	transactions []*Transaction
}

func NewTrManager() *TrManager {
	return &TrManager{make([]*Transaction, 0)}
}

func NewTrManagerFromOperations(operations []*Operation) *TrManager {
	trManager := NewTrManager()

	for _, op := range operations {
		trManager.AddTransaction(NewTransaction(op.ID()))
	}

	return trManager
}

func (trm *TrManager) AddTransaction(tr *Transaction) {
	if trm.Contains(tr.ID()) {
		return
	}

	trm.transactions = append(trm.transactions, tr)
}

func (trm *TrManager) Contains(id int) bool {
	for _, tr := range trm.transactions {
		if tr.ID() == id {
			return true
		}
	}

	return false
}

func (trm *TrManager) Transactions() []*Transaction {
	return trm.transactions
}

func (trm *TrManager) GetTransaction(id int) *Transaction {
	for idx := range trm.transactions {
		if trm.transactions[idx].ID() == id {
			return trm.transactions[idx]
		}
	}

	return nil
}

func (trm *TrManager) GetTransactionStatus(id int) (TrStatus, error) {
	for _, tr := range trm.transactions {
		if tr.ID() == id {
			return tr.Status(), nil
		}
	}

	return Finished, fmt.Errorf("não existe nenhuma transação na tabela com o id %d", id)
}

func (trm *TrManager) UpdateStatus(id int, status TrStatus) {
	for _, tr := range trm.transactions {
		if tr.ID() == id {
			tr.UpdateStatus(status)
		}
	}
}

func (trm *TrManager) PrintTransactions() {
	fmt.Println("|------------------------|")
	fmt.Println("|   TrID   |   Status    |")
	fmt.Println("|------------------------|")
	for _, transaction := range trm.transactions {
		fmt.Printf("|   %-6d |   %-9s |\n", transaction.ID(), transaction.Status())
	}
	fmt.Printf("|------------------------|\n\n")
}
