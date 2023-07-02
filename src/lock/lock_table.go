package lock

import (
	"fmt"
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

func (lt *LockTable) ReadLock(tr int, itemKey string) int {

	for _, lock := range lt.locks {
		if lock.LockType == WriteLock && lock.ItemKey == itemKey && lock.TrID != tr {
			return lock.TrID
		}
	}

	lock := NewLock(itemKey, tr, DefaultReadDuration, ReadLock)
	lt.locks = append(lt.locks, lock)
	return -1
}

func (lt *LockTable) WriteLock(tr int, itemKey string) int {

	for _, lock := range lt.locks {
		if lock.ItemKey == itemKey && lock.TrID != tr {
			return lock.TrID
		}

	}

	for _, lock := range lt.locks {
		// Upgrade de Lock
		if lock.TrID == tr && lock.LockType == ReadLock && lock.ItemKey == itemKey {
			lock.LockType = WriteLock
			lock.Duration = DefaultWriteDuration
			return -1
		}

	}

	lock := NewLock(itemKey, tr, DefaultReadDuration, WriteLock)
	lt.locks = append(lt.locks, lock)
	return -1
}

func (lt *LockTable) Unlock(tr int, itemKey string) {
	newLocks := make([]*Lock, 0)

	for idx := range lt.locks {
		if lt.locks[idx].TrID != tr || lt.locks[idx].ItemKey != itemKey {
			newLocks = append(newLocks, lt.locks[idx])
		}
	}

	lt.locks = newLocks
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
