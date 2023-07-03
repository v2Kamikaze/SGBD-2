package lock

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
	if !q.InQueue(item) {
		(*q) = append((*q), item)
	}
}

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *Queue) Dequeue() (int, error) {
	if len(*q) == 0 {
		return math.MinInt, errors.New("fila vazia")
	}

	e := (*q)[0]

	(*q) = (*q)[1:]

	return e, nil
}

func (q *Queue) RemoveAt(index int) error {
	if index < 0 || index >= len(*q) {
		return errors.New("índice fora dos limites")
	}

	*q = append((*q)[:index], (*q)[index+1:]...)
	return nil
}

func (q *Queue) First() (int, error) {
	if len(*q) == 0 {
		return math.MinInt, errors.New("fila vazia")
	}

	e := (*q)[0]

	return e, nil
}

func (q *Queue) Read() {
	fmt.Println(q)
}

func (q *Queue) InQueue(item int) bool {
	for _, v := range *q {
		if v == item {
			return true
		}
	}

	return false
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
