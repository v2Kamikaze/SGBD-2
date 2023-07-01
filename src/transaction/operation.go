package transaction

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
	id     int
	opType OperationType
	item   string
}

func NewOperation(id int, opType OperationType, item string) *Operation {
	return &Operation{id, opType, item}
}

func (op *Operation) ID() int {
	return op.id
}

func (op *Operation) Type() OperationType {
	return op.opType
}

func (op *Operation) Item() string {
	return op.item
}

func (op *Operation) String() string {
	if op.item == "" {
		return fmt.Sprintf("(ID: %d, OP_TYPE: %s)", op.id, op.opType)
	}

	return fmt.Sprintf("(ID: %d, OP_TYPE: %s, ITEM: %s)", op.id, op.opType, op.item)
}
