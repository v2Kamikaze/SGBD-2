package main

import (
	"github.com/v2Kamikaze/SGBD-2/src"
	"github.com/v2Kamikaze/SGBD-2/src/scheduler"
)

func main() {
	//q := "BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2)"
	q := "BT(1)r1(x)BT(2)w2(y)r1(y)w2(x)C(1)r2(z)C(2)"

	scheduling := src.ParseOperations(q)
	scheduler := scheduler.New(scheduling)
	scheduler.Start()
}

//  quando uma transação está esperando por um bloqueio em um objeto, ela não pode continuar fazendo operações em outros objetos
