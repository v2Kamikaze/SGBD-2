package lock

type TrScope int
type TrDuration int
type TrLockType int

const (
	Object TrScope = iota
	Predicate
)

const (
	Short TrDuration = iota
	Long
)

const (
	ReadLock TrLockType = iota
	WriteLock
)

type Lock struct {
	ItemKey  string
	TrID     int
	Scope    TrScope
	Duration TrDuration
	LockType TrLockType
}

func NewLock(key string, trID int, duration TrDuration, lockType TrLockType) *Lock {
	return &Lock{
		key,
		trID,
		Object,
		duration,
		lockType,
	}
}

func (scope TrScope) String() string {
	switch scope {
	case Object:
		return "Objeto"
	case Predicate:
		return "Predicado"
	default:
		return ""
	}
}

func (duration TrDuration) String() string {
	switch duration {
	case Short:
		return "Curta"
	case Long:
		return "Longa"
	default:
		return ""
	}
}

func (lockType TrLockType) String() string {
	switch lockType {
	case ReadLock:
		return "Leitura"
	case WriteLock:
		return "Escrita"
	default:
		return ""
	}
}
