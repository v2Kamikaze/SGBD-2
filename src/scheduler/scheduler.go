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
	//trDelayed  map[int]*lock.Queue
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

		fmt.Println(s.delayed)
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

		// A transação tá livre para executar, já que não está mais esperando
		if len(s.waitFor.GetNeighbors(op.ID())) == 0 && s.waitItem.IsWaiting(op.Item(), op.ID()) {
			s.trManager.UpdateStatus(op.ID(), transaction.Active)
			switch op.Type() {
			case transaction.ReadOp:
				if status, _ := s.trManager.GetTransactionStatus(op.ID()); status != transaction.Waiting {
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
								break // Não continua executando operações em outros objetos
							}
						} else {
							// Armazena as posições das operações conflitantes
							copy.Enqueue(idx)
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
							// Armazena as posições das operações conflitantes
							copy.Enqueue(idx)
							s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
							break // Não continua executando operações em outros objetos
						} else {
							fmt.Printf("Operação %d foi abortada\n", op.ID())
							// Transação atual é abortada (algoritmo Wait Die)
							s.trManager.UpdateStatus(trID, transaction.Aborted)
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

/*
func (s *Scheduler) Start2() {

	for idx, op := range s.scheduling {
		trID := op.ID()
		trStatus := s.trManager.GetTransaction(trID).Status()

		if trStatus != transaction.Aborted || op.Type() != transaction.BeginOp {
			switch op.Type() {
			case transaction.ReadOp:
				item := op.Item()
				if idTrInConflict := s.lockTable.ReadLock(trID, item); idTrInConflict != -1 {
					fmt.Printf("Houve conflito entre as operações da transação %d e %d\n", trID, idTrInConflict)
					s.waitFor.AddEdge(trID, idTrInConflict) // Cria uma aresta no grafo de conflitos (wait for)

					// Armazena as posições das operações conflitantes
					s.delayed.Enqueue(idx)

					if s.waitFor.HasCycle() {
						if trID < idTrInConflict {
							// Transação atual entra em espera (protocolo 2PL estrito)
							s.waitItem.EnqueueItem(item, trID)
							s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
						} else {
							fmt.Printf("Operação %d foi abortada\n", op.ID())

							// Transação atual é abortada (algoritmo Wait Die)
							s.trManager.UpdateStatus(trID, transaction.Aborted)
						}
					} else {
						// Armazena as posições das operações conflitantes
						s.delayed.Enqueue(idx)
					}
				}
			case transaction.WriteOp:
				item := op.Item()
				if idTrInConflict := s.lockTable.WriteLock(trID, item); idTrInConflict != -1 {
					fmt.Printf("Houve conflito entre as operações da transação %d e %d\n", trID, idTrInConflict)
					s.waitFor.AddEdge(trID, idTrInConflict) // Cria uma aresta no grafo de conflitos (wait for)
					// Armazena as posições das operações conflitantes
					s.delayed.Enqueue(idx)

					// Algoritmo Wait Die
					if s.waitFor.HasCycle() {
						if trID < idTrInConflict {
							// Armazena as posições das operações conflitantes
							s.delayed.Enqueue(idx)

							// Transação atual entra em espera (algoritmo Wait Die)
							s.waitItem.EnqueueItem(item, trID)
							s.trManager.UpdateStatus(op.ID(), transaction.Waiting)
						} else {
							fmt.Printf("Operação %d foi abortada\n", op.ID())
							// Transação atual é abortada (algoritmo Wait Die)
							s.trManager.UpdateStatus(trID, transaction.Aborted)
						}
					} else {
						// Transação atual obtém o bloqueio (2PL estrito)
						s.lockTable.WriteLock(trID, item)
					}
				}
			case transaction.CommitOp:
				for _, iop := range s.scheduling {
					// Liberando todos os locks da transação
					if iop.ID() == trID {
						s.lockTable.Unlock(trID, iop.Item())
					}
				}
				// Atualizando o status na TrManager
				s.trManager.UpdateStatus(trID, transaction.Finished)
				s.waitItem.Dequeue(op.Item())   // Remove a transação da fila de espera do item
				s.waitFor.RemoveVertex(op.ID()) // Remover do grafo a aresta que liga Ti até Tj
			}
		}

		// Processa as transações em espera
		s.ProcessWaitingTransactions()

		fmt.Println("Iteração ", idx)
		fmt.Println(s.delayed)
		s.ShowTables()

	}
}

func (s *Scheduler) ProcessWaitingTransactions() {
	for {
		idx, err := s.delayed.Dequeue()
		if err != nil {
			break // Fila vazia, sair do loop
		}

		op := s.scheduling[idx]
		trID := op.ID()

		if s.trManager.GetTransaction(trID).Status() == transaction.Aborted {
			continue // Transação abortada, passa para a próxima transação
		}

		// Verifica se é um bloqueio de leitura ou escrita
		isWriteLock := op.Type() == transaction.WriteOp

		// Tenta obter o bloqueio novamente
		var idTrInConflict int
		if isWriteLock {
			idTrInConflict = s.lockTable.WriteLock(trID, op.Item())
		} else {
			idTrInConflict = s.lockTable.ReadLock(trID, op.Item())
		}

		if idTrInConflict == -1 {
			// Transação obteve o bloqueio (2PL estrito)
			if s.delayed.IsEmpty() {
				s.waitItem.DeleteKey(op.Item()) // Remove o item da lista de espera se a fila estiver vazia
			}

			// Executa a operação
		} else {
			// Transação ainda em espera, adiciona de volta à fila
			s.delayed.Enqueue(idx)
			break
		}
	}
}

func (s *Scheduler) ProcessWaitingTransactions2() {
	for item, queue := range *s.waitItem {
		for {
			trID, err := queue.Dequeue()
			if err != nil {
				break // Fila vazia, passa para o próximo item
			}

			// Verifica se a transação está abortada
			if s.trManager.GetTransaction(trID).Status() == transaction.Aborted {
				continue // Transação abortada, passa para a próxima transação
			}

			// Verifica se é um bloqueio de leitura ou escrita
			isWriteLock := s.lockTable.IsWriteLock(item)

			// Tenta obter o bloqueio novamente
			var idTrInConflict int
			if isWriteLock {
				idTrInConflict = s.lockTable.WriteLock(trID, item)
			} else {
				idTrInConflict = s.lockTable.ReadLock(trID, item)
			}

			if idTrInConflict == -1 {
				// Transação obteve o bloqueio (2PL estrito)
				if len(*queue) == 0 {
					s.waitItem.DeleteKey(item) // Remove o item da lista de espera se a fila estiver vazia
				}
			} else {
				// Transação ainda em espera, adiciona de volta à fila
				queue.Enqueue(trID)
				break
			}
		}
	}
}
*/
