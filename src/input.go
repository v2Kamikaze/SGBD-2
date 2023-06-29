package src

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func InputTransactions() []string {

	scan := bufio.NewReader(os.Stdin)

	if input, err := scan.ReadBytes('\n'); err == nil {
		input = bytes.Replace(input, []byte("\r\n"), []byte(""), -1)
		for _, operation := range ParseOperations(string(input)) {
			fmt.Println("OP: ", operation)
		}
	}

	return nil
}

func ParseOperations(src string) []*Operation {
	ops := strings.SplitAfter(src, ")")
	ops = ops[:len(ops)-1]

	operations := make([]*Operation, len(ops))

	for i, op := range ops {
		operations[i] = ParseOperation(op)
	}

	return operations
}

func ParseOperation(op string) *Operation {
	var operation *Operation
	op = strings.Replace(op, "(", " ", -1)
	op = strings.Replace(op, ")", " ", -1)
	keywords := strings.Split(op, " ")
	opQuery := keywords[0]
	param := keywords[1]

	if strings.Contains(opQuery, "BT") {
		id := StrToInt(param)
		operation = NewOperation(id, OperationTypeFromStr("BT"), "")

	} else if strings.Contains(opQuery, "r") {
		id := GetTransactionID(opQuery, "r")
		operation = NewOperation(id, OperationTypeFromStr("r"), param)

	} else if strings.Contains(opQuery, "w") {
		id := GetTransactionID(opQuery, "w")
		operation = NewOperation(id, OperationTypeFromStr("w"), param)

	} else if strings.Contains(opQuery, "C") {
		id := StrToInt(param)
		operation = NewOperation(id, OperationTypeFromStr("C"), "")
	}

	return operation
}

func StrToInt(str string) int {
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("Não foi possível converter o valor %s em um valor inteiro. Erro: %+v", str, err)
	}
	return id
}

func GetTransactionID(opQuery, opType string) int {
	typeAndID := strings.SplitAfter(opQuery, opType)
	return StrToInt(typeAndID[1])
}

// BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2)
