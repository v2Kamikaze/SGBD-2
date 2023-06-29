package main

import (
	"github.com/v2Kamikaze/SGBD-2/src"
)

func main() {
	/* for _, operation := range src.ParseOperations("BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2)") {
		fmt.Printf("OP: %s\n", operation)
	} */

	wi := src.NewWaitItem()
	wi.EnqueueItem("x", 1)
	wi.EnqueueItem("x", 2)
	wi.EnqueueItem("z", 1)
	wi.EnqueueItem("y", 3)
	wi.EnqueueItem("y", 5)

	wi.ReadAll()

}
