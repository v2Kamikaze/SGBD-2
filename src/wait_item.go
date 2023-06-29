package src

import "fmt"

type WaitItem map[string]*Queue

func NewWaitItem() *WaitItem {
	wi := make(WaitItem)
	return &wi
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
	waitItem := *wi
	return waitItem[itemKey].Dequeue()
}

func (wi WaitItem) ReadAll() {
	for key, queue := range wi {
		fmt.Printf("%s : ", key)
		queue.Read()
	}
}
