package transaction

type OperationsTable map[int][]*Operation

func NewOperationsTable() OperationsTable {
	return make(OperationsTable)
}

func (opt *OperationsTable) AddOperation(operation *Operation) {
	table := *opt

	if _, ok := table[operation.ID()]; !ok {
		table[operation.ID()] = append(table[operation.ID()], operation)
		return
	}

	table[operation.ID()] = append(table[operation.ID()], operation)
}
