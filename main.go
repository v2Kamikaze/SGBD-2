package main

import (
	"fmt"

	"github.com/v2Kamikaze/SGBD-2/src"
	"github.com/v2Kamikaze/SGBD-2/src/scheduler"
)

func main() {
	fmt.Println("Escolha o nível de isolamento:")
	fmt.Println("1 - Read Uncommitted")
	fmt.Println("2 - Read Committed")
	fmt.Println("3 - Repeatable Read")
	fmt.Println("4 - Serializable")

	// var isolationLevel scheduler.IsolationLevel

	// for {
	// 	fmt.Print("Digite o número correspondente ao nível de isolamento: ")
	// 	var input string
	// 	fmt.Scanln(&input)
	// 	input = strings.TrimSpace(input)
	// 	if input == "1" {
	// 		isolationLevel = scheduler.ReadUncommitted
	// 		break
	// 	} else if input == "2" {
	// 		isolationLevel = scheduler.ReadCommitted
	// 		break
	// 	} else if input == "3" {
	// 		isolationLevel = scheduler.RepeatableRead
	// 		break
	// 	} else if input == "4" {
	// 		isolationLevel = scheduler.Serializable
	// 		break
	// 	} else {
	// 		fmt.Println("Opção inválida. Digite novamente.")
	// 	}
	// }

	fmt.Print("Digite a sequência de operações de escalonamento (exemplo: BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2)): ")
	var transactionString string
	fmt.Scanln(&transactionString)

	scheduling := src.ParseOperations(transactionString)

	scheduler := scheduler.New(scheduling)
	// scheduler.SetIsolationLevel(isolationLevel)

	fmt.Println("Executando escalonamento...")
	scheduler.Start()
}

//  quando uma transação está esperando por um bloqueio em um objeto, ela não pode continuar fazendo operações em outros objetos
