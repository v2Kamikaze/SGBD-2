package main

import (
	"github.com/v2Kamikaze/SGBD-2/src"
)

func main() {
	/* for _, operation := range src.ParseOperations("BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2)") {
		fmt.Printf("OP: %s\n", operation)
	} */

	graph := src.NewGraph()
	graph.AddEdge(1, 2)
	graph.AddEdge(3, 2)
	graph.PrintGraphTable()

	graph.AddEdge(2, 3)

	if graph.HasCycle() {
		graph.RemoveEdge(2, 3)
	}

	graph.PrintGraphTable()

}
