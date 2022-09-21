package utils

const WORKER_PORT int = 9999
const MASTER_PORT int = 9001
const WORKER_IP string = "localhost"

type workerType int

const (
	Mapper      workerType = 1
	Reducer     workerType = 2
	NotAssigned workerType = 0
)

type JoinRequest struct {
	IP   string
	Type workerType
}

type WorkerInfo struct {
	IP string
}
