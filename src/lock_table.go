package src

type TrScope int
type TrDuration int
type TrLockType int

const (
	object TrScope = iota
	predicate
)

const (
	short TrDuration = iota
	long
)

const (
	read TrLockType = iota
	write
)

type LockTable struct {
	IDItem   int
	TrID     int
	Scope    TrScope
	Duration TrDuration
	LockType TrLockType
}

func NewLockTable() *LockTable {
	return &LockTable{}
}

func (lt *LockTable) RL(tr *TrManager, d int) {

}

func (lt *LockTable) WL(tr *TrManager, d int) {

}

func (lt *LockTable) UL(tr *TrManager, d int) {

}
