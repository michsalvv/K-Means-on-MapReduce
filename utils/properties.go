package utils

type workerType int

const (
	MapperType  workerType = 1
	ReducerType workerType = 2
)

type JoinRequest struct {
	Port int
	Type workerType
}
