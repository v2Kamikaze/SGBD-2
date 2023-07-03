package lock

import (
	"errors"
	"fmt"
)

type WaitItem map[string]*Queue

func NewWaitItem() *WaitItem {
	wi := make(WaitItem)
	return &wi
}

func (wi *WaitItem) Self() map[string]*Queue {
	return *wi
}

func (wi *WaitItem) EnqueueItem(itemKey string, trID int) {
	waitItem := *wi

	if _, ok := waitItem[itemKey]; ok {
		waitItem[itemKey].Enqueue(trID)
	} else {
		waitItem[itemKey] = NewQueue()
		waitItem[itemKey].Enqueue(trID)
	}
}

func (wi *WaitItem) Dequeue(itemKey string) (int, error) {
	queue, ok := (*wi)[itemKey]
	if !ok {
		return 0, errors.New("fila vazia")
	}

	trID, err := queue.Dequeue()
	if err != nil {
		// Fila vazia, remove o item do WaitItem
		delete(*wi, itemKey)
	}

	return trID, err
}

func (wi *WaitItem) InQueue(id int) bool {
	for key := range *wi {
		if (*wi)[key].InQueue(id) {
			return true
		}
	}

	return false
}

func (wi *WaitItem) DeleteKey(itemKey string) {
	delete(*wi, itemKey)
}

func (wi *WaitItem) IsWaiting(itemKey string, id int) bool {
	if _, ok := (*wi)[itemKey]; ok {
		return (*wi)[itemKey].InQueue(id)
	}

	return false
}

func (wi WaitItem) PrintWaitList() {
	fmt.Println("|------------- Wait Item -------------|")
	for key, queue := range wi {
		fmt.Printf("%s : ", key)
		queue.Read()
	}
	fmt.Printf("|-------------------------------------|\n\n")
}
