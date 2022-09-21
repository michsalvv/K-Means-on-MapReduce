package utils

import "strings"

const WORKER_PORT int = 9999
const MASTER_PORT int = 9001
const WORKER_IP string = "localhost"

type WorkerType int

const (
	Mapper  WorkerType = 1
	Reducer WorkerType = 2
	NONE    WorkerType = -1
)

func DetectTaskType(workerType string) WorkerType {
	if strings.Compare(workerType, "reducer") == 0 {
		return WorkerType(2)
	}
	if strings.Compare(workerType, "mapper") == 0 {
		return WorkerType(1)
	}
	return WorkerType(-1)
}

type JoinRequest struct {
	IP   string
	Type WorkerType
}

type WorkerInfo struct {
	IP string
}
