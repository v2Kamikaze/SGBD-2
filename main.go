package main

import (
	"github.com/v2Kamikaze/SGBD-2/src"
	"github.com/v2Kamikaze/SGBD-2/src/scheduler"
)

func main() {

	scheduling := src.ParseOperations("BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2)")

	scheduler := scheduler.New(scheduling)
	scheduler.Start()
}
