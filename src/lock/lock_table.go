package lock

import (
	"fmt"

	"github.com/v2Kamikaze/SGBD-2/src/transaction"
)

var DefaultReadDuration TrDuration = Short
var DefaultWriteDuration TrDuration = Long

type LockTable struct {
	locks []*Lock
}

func NewLockTable() *LockTable {
	return &LockTable{make([]*Lock, 0)}
}

func (lt *LockTable) ReadLock(tr *transaction.Transaction, itemKey string) int {

	if lt.CheckLockExistence(tr.ID(), ReadLock, itemKey) {
		return -1
	}

	for _, lock := range lt.locks {
		if lock.LockType == WriteLock && lock.ItemKey == itemKey && lock.TrID != tr.ID() {
			return lock.TrID
		}
	}

	lock := NewLock(itemKey, tr.ID(), DefaultReadDuration, ReadLock)
	lt.locks = append(lt.locks, lock)
	return -1
}

func (lt *LockTable) WriteLock(tr *transaction.Transaction, itemKey string) int {

	if lt.CheckLockExistence(tr.ID(), WriteLock, itemKey) {
		return -1
	}

	for _, lock := range lt.locks {
		if lock.ItemKey == itemKey && lock.TrID != tr.ID() {
			return lock.TrID
		}

	}

	lock := NewLock(itemKey, tr.ID(), DefaultReadDuration, WriteLock)
	lt.locks = append(lt.locks, lock)
	return -1
}

func (lt *LockTable) Unlock(tr *transaction.Transaction, itemKey string) {
	newLocks := make([]*Lock, 0)

	for idx := range lt.locks {
		if lt.locks[idx].TrID != tr.ID() || lt.locks[idx].ItemKey != itemKey {
			newLocks = append(newLocks, lt.locks[idx])
		}
	}

	lt.locks = newLocks
}

func (lt *LockTable) IsWriteLock(itemKey string) bool {
	for _, lock := range lt.locks {
		if lock.ItemKey == itemKey && lock.LockType == WriteLock {
			return true
		}
	}
	return false
}

func (lt *LockTable) CheckLockExistence(trID int, lockType TrLockType, itemKey string) bool {
	for _, lock := range lt.locks {
		if lock.TrID == trID && lock.LockType == lockType && lock.ItemKey == itemKey {
			return true
		}
	}
	return false
}

func (lt *LockTable) PrintTable() {
	fmt.Println("|-------------------------------------------------------------|")
	fmt.Println("|   IDItem   |   TrID   |   Escopo   |   Duração   |   Tipo   |")
	fmt.Println("|-------------------------------------------------------------|")
	for _, lock := range lt.locks {
		fmt.Printf("|  %-9s | %-8d | %-10s | %-11s | %-8s |\n", lock.ItemKey, lock.TrID, lock.Scope, lock.Duration, lock.LockType)
	}
	fmt.Printf("|-------------------------------------------------------------|\n\n")
}
