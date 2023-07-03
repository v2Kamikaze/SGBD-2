package scheduler

import (
	"fmt"

	"github.com/v2Kamikaze/SGBD-2/src/lock"
	"github.com/v2Kamikaze/SGBD-2/src/transaction"
)

type Scheduler struct {
	waitFor   *lock.Graph
	waitItem  *lock.WaitItem
	lockTable *lock.LockTable
	trManager *transaction.TrManager

	scheduling []*transaction.Operation
	delayed    *lock.Queue
}

func New(scheduling []*transaction.Operation) *Scheduler {
	return &Scheduler{
		waitFor:    lock.NewGraph(),
		waitItem:   lock.NewWaitItem(),
		lockTable:  lock.NewLockTable(),
		trManager:  transaction.NewTrManagerFromOperations(scheduling),
		scheduling: scheduling,
		delayed:    lock.NewQueue(),
		//trDelayed:  make(map[int]*lock.Queue),
	}
}

func (s *Scheduler) ShowTables() {
	s.trManager.PrintTransactions()
	s.lockTable.PrintTable()
	s.waitFor.PrintGraphTable()
	s.waitItem.PrintWaitList()
}

func (s *Scheduler) Start() {

	for idx, op := range s.scheduling {
		fmt.Println("Iteração ", idx)

		trID := op.ID()
		trStatus := s.trManager.GetTransaction(trID).Status()

		if trStatus == transaction.Aborted {
			continue
		}

		if trStatus != transaction.Aborted || op.Type() != transaction.BeginOp {
			switch op.Type() {
			case transaction.ReadOp:

				item := op.Item()
				tr := s.trManager.GetTransaction(op.ID())

				if tr.Status() == transaction.Waiting {
					fmt.Printf("Operação %s esperando.\n", op)
					s.delayed.Enqueue(idx)
					break
				}

				if idTrInConflict := s.lockTable.ReadLock(tr, item); idTrInConflict != -1 {

					fmt.Printf("Houve conflito entre as operações da transação %d e %d\n", trID, idTrInConflict)
					s.waitFor.AddEdge(trID, idTrInConflict) // Cria uma aresta no grafo de conflitos (wait for)
					s.waitItem.EnqueueItem(item, trID)

					// Armazena as posições das operações conflitantes
					s.delayed.Enqueue(idx)
					s.trManager.UpdateStatus(op.ID(), transaction.Waiting)

					if s.waitFor.HasCycle() {
						if trID < idTrInConflict {
							s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
							break // Não continua executando operações em outros objetos
						} else {
							fmt.Printf("Operação %d foi abortada\n", op.ID())
							// Transação atual é abortada (algoritmo Wait Die)
							s.trManager.UpdateStatus(trID, transaction.Aborted)
							s.waitFor.RemoveVertex(op.ID())
							s.waitItem.Dequeue(op.Item())
							for _, iop := range s.scheduling {
								tr := s.trManager.GetTransaction(op.ID())
								// Liberando todos os locks da transação
								if iop.ID() == trID {
									s.lockTable.Unlock(tr, iop.Item())
								}
							}
							break // Não continua executando operações em outros objetos
						}
					} else {
						// Armazena as posições das operações conflitantes
						s.delayed.Enqueue(idx)
					}
				}

			case transaction.WriteOp:
				item := op.Item()
				tr := s.trManager.GetTransaction(op.ID())

				if tr.Status() == transaction.Waiting {
					fmt.Printf("Operação %s esperando.\n", op)
					s.delayed.Enqueue(idx)
					break
				}

				if idTrInConflict := s.lockTable.WriteLock(tr, item); idTrInConflict != -1 {

					fmt.Printf("Houve conflito entre as operações da transação %d e %d\n", trID, idTrInConflict)
					s.waitFor.AddEdge(trID, idTrInConflict) // Cria uma aresta no grafo de conflitos (wait for)
					s.waitItem.EnqueueItem(item, trID)

					// Armazena as posições das operações conflitantes
					s.delayed.Enqueue(idx)
					s.trManager.UpdateStatus(op.ID(), transaction.Waiting)

					// Algoritmo Wait Die
					if s.waitFor.HasCycle() {
						if trID < idTrInConflict {
							// Armazena as posições das operações conflitantes
							s.delayed.Enqueue(idx)
							s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
							break // Não continua executando operações em outros objetos
						} else {
							fmt.Printf("Operação %d foi abortada\n", op.ID())
							// Transação atual é abortada (algoritmo Wait Die)
							s.trManager.UpdateStatus(trID, transaction.Aborted)
							s.waitFor.RemoveVertex(op.ID())
							s.waitItem.Dequeue(op.Item())
							for _, iop := range s.scheduling {
								tr := s.trManager.GetTransaction(op.ID())
								// Liberando todos os locks da transação
								if iop.ID() == trID {
									s.lockTable.Unlock(tr, iop.Item())
								}
							}
							break // Não continua executando operações em outros objetos
						}
					} else {
						// Transação atual obtém o bloqueio (2PL estrito)
						s.lockTable.WriteLock(tr, item)
					}
				}

			case transaction.CommitOp:
				// Não pode ocorrer o commit se ainda estiver esperando
				if !s.waitItem.InQueue(op.ID()) {
					for _, iop := range s.scheduling {
						tr := s.trManager.GetTransaction(op.ID())
						// Liberando todos os locks da transação
						if iop.ID() == trID {
							s.lockTable.Unlock(tr, iop.Item())
						}
					}
					// Atualizando o status na TrManager
					s.trManager.UpdateStatus(trID, transaction.Finished)
					s.waitItem.Dequeue(op.Item())   // Remove a transação da fila de espera do item
					s.waitFor.RemoveVertex(op.ID()) // Remover do grafo a aresta que liga Ti até Tj

				} else {
					s.delayed.Enqueue(idx)
					s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
				}
			}
		}

		s.TryRestartWaitingTransactions()

		fmt.Println("Delayed: ", s.delayed)
		s.ShowTables()
	}
}

func (s *Scheduler) TryRestartWaitingTransactions() {
	if s.delayed.IsEmpty() {
		return
	}

	copy := (*(s.delayed))[:]

	for _, idx := range *(s.delayed) {
		op := s.scheduling[idx]
		trID := op.ID()
		trStatus, _ := s.trManager.GetTransactionStatus(trID)

		if trStatus == transaction.Aborted {
			continue
		}

		// A transação tá livre para executar, já que não está mais esperando
		if len(s.waitFor.GetNeighbors(op.ID())) == 0 {

			s.trManager.UpdateStatus(op.ID(), transaction.Active)
			switch op.Type() {
			case transaction.ReadOp:
				item := op.Item()
				tr := s.trManager.GetTransaction(op.ID())

				if tr.Status() == transaction.Waiting {
					fmt.Printf("Operação %s esperando.\n", op)
					copy.Enqueue(idx)
					break
				}

				if idTrInConflict := s.lockTable.ReadLock(tr, item); idTrInConflict != -1 {

					fmt.Printf("Houve conflito entre as operações da transação %d e %d\n", trID, idTrInConflict)
					s.waitFor.AddEdge(trID, idTrInConflict) // Cria uma aresta no grafo de conflitos (wait for)
					s.waitItem.EnqueueItem(item, trID)

					// Armazena as posições das operações conflitantes
					copy.Enqueue(idx)
					s.trManager.UpdateStatus(op.ID(), transaction.Waiting)

					if s.waitFor.HasCycle() {
						if trID < idTrInConflict {
							s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
							break // Não continua executando operações em outros objetos
						} else {
							fmt.Printf("Operação %d foi abortada\n", op.ID())

							// Transação atual é abortada (algoritmo Wait Die)
							s.trManager.UpdateStatus(trID, transaction.Aborted)
							s.waitFor.RemoveVertex(op.ID())
							s.waitItem.Dequeue(op.Item())
							for _, iop := range s.scheduling {
								tr := s.trManager.GetTransaction(op.ID())
								// Liberando todos os locks da transação
								if iop.ID() == trID {
									s.lockTable.Unlock(tr, iop.Item())
								}
							}

							break // Não continua executando operações em outros objetos
						}
					}
				} else {
					s.waitItem.Dequeue(op.Item())
					copy.Dequeue()
				}

			case transaction.WriteOp:
				item := op.Item()
				tr := s.trManager.GetTransaction(op.ID())

				if tr.Status() == transaction.Waiting {
					fmt.Printf("Operação %s esperando.\n", op)
					copy.Enqueue(idx)
					break
				}

				if idTrInConflict := s.lockTable.WriteLock(tr, item); idTrInConflict != -1 {

					fmt.Printf("Houve conflito entre as operações da transação %d e %d\n", trID, idTrInConflict)
					s.waitFor.AddEdge(trID, idTrInConflict) // Cria uma aresta no grafo de conflitos (wait for)
					s.waitItem.EnqueueItem(item, trID)

					// Armazena as posições das operações conflitantes
					copy.Enqueue(idx)
					s.trManager.UpdateStatus(op.ID(), transaction.Waiting)

					// Algoritmo Wait Die
					if s.waitFor.HasCycle() {
						if trID < idTrInConflict {
							s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
							break // Não continua executando operações em outros objetos
						} else {
							fmt.Printf("Operação %d foi abortada\n", op.ID())
							// Transação atual é abortada (algoritmo Wait Die)
							s.trManager.UpdateStatus(trID, transaction.Aborted)
							s.waitFor.RemoveVertex(op.ID())
							s.waitItem.Dequeue(op.Item())
							for _, iop := range s.scheduling {
								tr := s.trManager.GetTransaction(op.ID())
								// Liberando todos os locks da transação
								if iop.ID() == trID {
									s.lockTable.Unlock(tr, iop.Item())
								}
							}
							break // Não continua executando operações em outros objetos
						}
					} else {
						// Transação atual obtém o bloqueio (2PL estrito)
						s.lockTable.WriteLock(tr, item)
					}
				} else {
					s.waitItem.Dequeue(op.Item())
					copy.Dequeue()
				}

			case transaction.CommitOp:
				// Não pode ocorrer o commit se ainda estiver esperando
				if !s.waitItem.InQueue(op.ID()) {
					for _, iop := range s.scheduling {
						tr := s.trManager.GetTransaction(op.ID())
						// Liberando todos os locks da transação
						if iop.ID() == trID {
							s.lockTable.Unlock(tr, iop.Item())
						}
					}
					// Atualizando o status na TrManager
					s.trManager.UpdateStatus(trID, transaction.Finished)
					s.waitItem.Dequeue(op.Item())   // Remove a transação da fila de espera do item
					s.waitFor.RemoveVertex(op.ID()) // Remover do grafo a aresta que liga Ti até Tj
				} else {
					copy.Enqueue(idx)
					s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
				}
			}
		}
	}

	s.delayed = &copy

}
