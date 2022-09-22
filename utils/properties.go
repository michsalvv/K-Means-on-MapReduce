package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

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

func Wait() {

	log.Println("Press Enter...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func ViewClusters(clusters [][]Point, numClusters int, printCluster bool) {
	for i := 0; i < numClusters; i++ {
		log.Printf("Mapped [%d] points on cluster [#%d]", len(clusters[i]), i)
		if printCluster {
			log.Print(clusters[i])
		}
	}
}

type JoinRequest struct {
	IP   string
	Type WorkerType
}

type WorkerInfo struct {
	IP string
}

type InputKMeans struct {
	Dataset  string
	Clusters int
}

type Point struct {
	Values []float64 //d-dimensional
}
type MapperInput struct {
	Chunk     []Point
	Centroids []Point
}

type MapperResponse struct {
	Clusters [][]Point //[ClusterIndex][P1, P2, P3]
	IP       string
}

type ReducerResponse struct {
	Centroids []Point
	IP        string
}
