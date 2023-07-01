package lock

import (
	"fmt"

	"github.com/v2Kamikaze/SGBD-2/src/transaction"
)

type TrScope int
type TrDuration int
type TrLockType int

var DefaultReadDuration TrDuration = Short
var DefaultWriteDuration TrDuration = Long

type LockTable struct {
	locks []*Lock
}

func NewLockTable() *LockTable {
	return &LockTable{make([]*Lock, 0)}
}

func (lt *LockTable) ReadLock(tr *transaction.Transaction, itemKey string) error {
	lock := NewLock(itemKey, tr.ID(), DefaultReadDuration, ReadLock)
	lt.locks = append(lt.locks, lock)

	return nil
}

func (lt *LockTable) WriteLock(tr *transaction.Transaction, itemKey string) error {
	lock := NewLock(itemKey, tr.ID(), DefaultReadDuration, WriteLock)
	lt.locks = append(lt.locks, lock)
	return nil

}

func (lt *LockTable) Unlock(tr *transaction.Transaction, itemKey string) error {

	return nil

}

func (lt *LockTable) PrintTable() {
	fmt.Println("|-------------------------------------------------------------|")
	fmt.Println("|   IDItem   |   TrID   |   Escopo   |   Duração   |   Tipo   |")
	fmt.Println("|-------------------------------------------------------------|")
	for _, lock := range lt.locks {
		fmt.Printf("|  %-9s | %-8d | %-10s | %-11s | %-8s |\n", lock.ItemKey, lock.TrID, lock.Scope, lock.Duration, lock.LockType)
	}
	fmt.Println("|-------------------------------------------------------------|")
}
