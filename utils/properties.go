package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// const WORKER_PORT int = 9999
// const MASTER_PORT int = 9001
// const WORKER_IP string = "localhost"
const DATASET_DIR string = "/go/src/kmeans-MR/datasets/"

// const TEST_FILE = "results.csv"
// const CONV_THRESH float64 = 0.001
// const COMBINER bool = false

type WorkerType int

const (
	Mapper  WorkerType = 1
	Reducer WorkerType = 2
	NONE    WorkerType = -1
)

const (
	NO_RES_ERROR      string = "not enough resources to perform the algorithm"
	NO_REDUCERS_ERROR string = "not enough reducers available in the cluster to perform the algorithm"
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

func GetConfiguration() Config {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal("Configuration file {config.yml} not found")
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("Error parsing configuration file {config.yml}")
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal("Error parsing configuration file {config.yml}")
	}

	return cfg
}

type Config struct {
	Server struct {
		MASTER_PORT string `yaml:"master_port", envconfig:"MASTER_PORT"`
		WORKER_PORT string `yaml:"worker_port", envconfig:"WORKE_PORT"`
		HOST        string `yaml:"host", envconfig:"SERVER_HOST"`
		DATASET_DIR string `yaml:"dataset_dir", envconfig:"DATASET_DIR"`
	} `yaml:"server"`
	Parameters struct {
		TEST_FILE_NAME string  `yaml:"test_file_name", envconfig:"TEST_FILE_NAME"`
		CONV_THRESH    float64 `yaml:"conv_thresh", envconfig:"CONV_THRESH"`
		COMBINER       bool    `yaml:"combiner", envconfig:"COMBINER"`
	} `yaml:"parameters"`
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
	Cluster []Point
	IP      string
}

type ReducerInput struct {
	Mappers    []WorkerInfo
	ClusterKey int
}
type ReducerResponse struct {
	Centroid Point
	IP       string
}

type Result struct {
	Centroids     []Point
	Iterations    int
	Error         int
	ExecutionTime time.Duration
}

type Triple struct {
	P        Point
	Distance float64
	Centroid Point
}
