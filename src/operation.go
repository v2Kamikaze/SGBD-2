package src

import "fmt"

type OperationType int

const (
	BeginOp OperationType = iota
	ReadOp
	WriteOp
	CommitOp
)

func OperationTypeFromStr(op string) OperationType {
	switch op {
	case "BT":
		return BeginOp
	case "r":
		return ReadOp
	case "w":
		return WriteOp
	case "C":
		return CommitOp
	default:
		panic("Operação inválida: " + op)
	}
}

func (opt OperationType) String() string {
	switch opt {
	case BeginOp:
		return "BT"
	case ReadOp:
		return "r"
	case WriteOp:
		return "w"
	default:
		return "C"
	}
}

type Operation struct {
	ID   int
	Type OperationType
	Item string
}

func NewOperation(id int, opType OperationType, item string) *Operation {
	return &Operation{id, opType, item}
}

func (op *Operation) String() string {
	if op.Item == "" {
		return fmt.Sprintf("(ID: %d, OP_TYPE: %s)", op.ID, op.Type)
	}

	return fmt.Sprintf("(ID: %d, OP_TYPE: %s, ITEM: %s)", op.ID, op.Type, op.Item)
}
