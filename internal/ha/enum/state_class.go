package enum

type StateClass string

const (
	StateClassMeasurement     StateClass = "measurement"
	StateClassTotal           StateClass = "total"
	StateClassTotalIncreasing StateClass = "total_increasing"
)
