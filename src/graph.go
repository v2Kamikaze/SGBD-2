package src

import "fmt"

type Graph struct {
	vertices map[int][]int
}

func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[int][]int),
	}
}

func (g *Graph) AddEdge(src, dest int) {
	g.vertices[src] = append(g.vertices[src], dest)
}

func (g *Graph) RemoveEdge(src, dest int) {
	neighbors := g.vertices[src]
	for i, neighbor := range neighbors {
		if neighbor == dest {
			g.vertices[src] = append(neighbors[:i], neighbors[i+1:]...)
			break
		}
	}
}

func (g *Graph) GetNeighbors(vertex int) []int {
	return g.vertices[vertex]
}

func (g *Graph) HasCycle() bool {
	visited := make(map[int]bool)
	recursionStack := make(map[int]bool)

	for vertex := range g.vertices {
		if isCyclic(vertex, visited, recursionStack, g) {
			return true
		}
	}

	return false
}

func isCyclic(vertex int, visited, recursionStack map[int]bool, graph *Graph) bool {
	visited[vertex] = true
	recursionStack[vertex] = true

	for _, neighbor := range graph.vertices[vertex] {
		if !visited[neighbor] {
			if isCyclic(neighbor, visited, recursionStack, graph) {
				return true
			}
		} else if recursionStack[neighbor] {
			return true
		}
	}

	recursionStack[vertex] = false
	return false
}

func (g *Graph) PrintGraphTable() {
	fmt.Println("Graph Table:")
	for vertex, neighbors := range g.vertices {
		fmt.Printf("%d -> ", vertex)
		for _, neighbor := range neighbors {
			fmt.Printf("%d ", neighbor)
		}
		fmt.Println()
	}
}
