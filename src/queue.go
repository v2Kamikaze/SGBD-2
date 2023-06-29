package src

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Queue []int

func NewQueue() *Queue {
	queue := make(Queue, 0)
	return &queue
}

func (q *Queue) Enqueue(item int) {
	(*q) = append((*q), item)
}

func (q *Queue) Dequeue() (int, error) {
	if len(*q) == 0 {
		return math.MinInt, errors.New("fila vazia")
	}

	e := (*q)[0]

	(*q) = (*q)[1:]

	return e, nil
}

func (q *Queue) Read() {
	fmt.Println(q)
}

func (q *Queue) String() string {

	lastIdx := len(*q) - 1

	str := "[ "

	for idx, item := range *q {
		if idx == lastIdx {
			str += strconv.Itoa(item)
		} else {
			str += strconv.Itoa(item) + " <- "
		}
	}

	str += " ]"

	return str
}
