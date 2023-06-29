package src

type TrStatus int

const (
	active TrStatus = iota
	finished
	aborted
	waiting
)

type TrManager struct {
	trID   int
	status TrStatus
}

func NewTrManager(id int) *TrManager {
	return &TrManager{id, active}
}

func (tm *TrManager) Status() TrStatus {
	return tm.status
}

func (tm *TrManager) UpdateStatus(status TrStatus) {
	tm.status = status
}
