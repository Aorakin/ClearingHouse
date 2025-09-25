package enum

const (
	StatusCreated  = "created"
	StatusRunning  = "running"
	StatusCanceled = "canceled"
	StatusStopped  = "stopped"
)

var (
	UsingStatuses = []string{StatusCreated, StatusRunning}
)

var (
	FinishedStatuses = []string{StatusCanceled, StatusStopped}
)
