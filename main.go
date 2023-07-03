package main

import (
	"github.com/v2Kamikaze/SGBD-2/src"
	"github.com/v2Kamikaze/SGBD-2/src/scheduler"
)

func main() {

	scheduling := src.ParseOperations("BT(1)r1(x)BT(2)w2(x)C(2)r2(y)r1(y)C(1)r2(z)")
	scheduler := scheduler.New(scheduling)
	scheduler.Start()
}

// rl1(x)r1(x)            rl2(y)r2(y)rl1(y)r1(y)C1rl2(z)r2(z)wl2(x)w2(x)C2
//			  wl2(x)w2(x)

//  quando uma transação está esperando por um bloqueio em um objeto, ela não pode continuar fazendo operações em outros objetos
