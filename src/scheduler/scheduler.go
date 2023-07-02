package scheduler

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/v2Kamikaze/SGBD-2/src/lock"
	"github.com/v2Kamikaze/SGBD-2/src/transaction"
)

type Scheduler struct {
	waitFor   *lock.Graph
	waitItem  *lock.WaitItem
	lockTable *lock.LockTable
	trManager *transaction.TrManager

	scheduling []*transaction.Operation
}

func New(scheduling []*transaction.Operation) *Scheduler {
	return &Scheduler{
		waitFor:    lock.NewGraph(),
		waitItem:   lock.NewWaitItem(),
		lockTable:  lock.NewLockTable(),
		trManager:  transaction.NewTrManagerFromOperations(scheduling),
		scheduling: scheduling,
	}
}

func (s *Scheduler) ShowTables() {
	s.trManager.PrintTransactions()
	s.lockTable.PrintTable()
	s.waitFor.PrintGraphTable()
	s.waitItem.PrintWaitList()
}

func (s *Scheduler) Start() {
	for _, op := range s.scheduling {

		fmt.Println("Chegou a operação: ", op)

		if s.trManager.GetTransaction(op.ID()).Status() != transaction.Aborted || op.Type() != transaction.BeginOp {

			switch op.Type() {
			// "BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2)"
			case transaction.ReadOp:
				if idTrInConflict := s.lockTable.ReadLock(op.ID(), op.Item()); idTrInConflict != -1 {

					fmt.Printf("Houve conflito entre as operaçãos da transação %d e %d\n", op.ID(), idTrInConflict)
					// Se houver conflito, criamos uma aresta no grafo de conflitos (wait for)
					s.waitFor.AddEdge(op.ID(), idTrInConflict)

					// Wait Die
					if s.waitFor.HasCycle() {
						if op.ID() < idTrInConflict {
							// Adicionamos a transação atual na lista de espera
							s.waitItem.EnqueueItem(op.Item(), op.ID())
						} else {
							s.trManager.UpdateStatus(op.ID(), transaction.Aborted)
						}

					}

				}
				// "BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2)"
			case transaction.WriteOp:
				if idTrInConflict := s.lockTable.WriteLock(op.ID(), op.Item()); idTrInConflict != -1 {
					fmt.Printf("Houve conflito entre as operaçãos da transação %d e %d\n", op.ID(), idTrInConflict)
					// Se houver conflito, criamos uma aresta no grafo de conflitos (wait for)
					s.waitFor.AddEdge(op.ID(), idTrInConflict)

					// Adicionamos a transação atual na lista de espera
					s.waitItem.EnqueueItem(op.Item(), op.ID())

					// Wait Die
					if s.waitFor.HasCycle() {
						if op.ID() < idTrInConflict {
							// Adicionamos a transação atual na lista de espera
							s.waitItem.EnqueueItem(op.Item(), op.ID())
						} else {
							s.trManager.UpdateStatus(op.ID(), transaction.Aborted)
						}

					}
				}
			case transaction.CommitOp:
				for _, iop := range s.scheduling {
					// Liberando todos os locks da transação
					if iop.ID() == op.ID() {
						s.lockTable.Unlock(op.ID(), iop.Item())
					}
				}

				// Atualizando o status na TrManager
				s.trManager.UpdateStatus(op.ID(), transaction.Finished)

			}
		}

		time.Sleep(time.Second * 2)
		s.ShowTables()

	}

}

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
