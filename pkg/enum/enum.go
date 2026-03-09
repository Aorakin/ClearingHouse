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

// resourceTypeOrder defines the canonical display order: CPU=0, GPU=1, RAM=2, others=3.
var resourceTypeOrder = map[string]int{
	"CPU": 0,
	"GPU": 1,
	"RAM": 2,
}

func ResourceTypeOrder(name string) int {
	if v, ok := resourceTypeOrder[name]; ok {
		return v
	}
	return len(resourceTypeOrder)
}
