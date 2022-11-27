package main

import (
	"bufio"
	"fmt"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

// const masterPath string = "server/master/"

var mappers []utils.WorkerInfo
var reducers []utils.WorkerInfo

type Master int

// RPC exposed to workers to register for the service
func (m *Master) JoinMR(req utils.JoinRequest, reply *int) error {

	log.Printf("Request received from [%s] type [%d]", req.IP, req.Type)
	switch req.Type {

	case utils.Mapper:
		*reply, mappers = addWorker(req, mappers)

	case utils.Reducer:
		*reply, reducers = addWorker(req, reducers)
	}

	return nil
}

// add Mapper/Worker Information to Mapper/Reducer list
func addWorker(req utils.JoinRequest, workersList []utils.WorkerInfo) (int, []utils.WorkerInfo) {
	for _, x := range workersList {
		if x.IP == req.IP {
			log.Printf("Request from [%s] declined", req.IP)
			return -1, workersList
		}
	}
	log.Printf("Worker [%s] accepted in cluster\n", req.IP)
	workersList = append(workersList, utils.WorkerInfo{IP: req.IP})
	return 0, workersList
}

// remove add Mapper/Worker Information
func removeWorker(req utils.JoinRequest, workersList []utils.WorkerInfo) (int, []utils.WorkerInfo) {
	var toRemove int
	for index, x := range workersList {
		if x.IP == req.IP {
			toRemove = index
			break
		}
	}
	workersList = append(workersList[:toRemove], workersList[toRemove+1:]...)
	return 0, workersList
}

func (m *Master) ExitMR(req utils.JoinRequest, reply *int) error {
	log.Printf("Worker [%s] disconnected", req.IP)
	switch req.Type {

	case utils.Mapper:
		*reply, mappers = removeWorker(req, mappers)

	case utils.Reducer:
		*reply, reducers = removeWorker(req, reducers)
	}
	return nil
}

// RPC exposed to user client for start kmeans Algorithm
func (m *Master) KMeans(in utils.InputKMeans, reply *utils.Result) error {

	clusterError := checkAvailability(in, mappers, reducers)

	if clusterError != nil {
		return clusterError
	}
	fmt.Print("\n\n")
	log.Print("--------------------------------------------------------")
	log.Printf("Starting K-Means Clustering")
	log.Printf("Dataset: {%s}", in.Dataset)

	start := time.Now()
	datasetPath := utils.DATASET_DIR + strings.Replace(in.Dataset, ".csv", "/", 1) + in.Dataset

	file, err := os.Open(datasetPath)
	if err != nil {
		log.Print(err.Error())
		return nil
	}

	reader := bufio.NewReader(file)
	var points []utils.Point
	point, err := readPoint(reader)
	for err == nil {
		if point.Values == nil {
			break
		}
		points = append(points, point)
		point, err = readPoint(reader)
	}

	log.Printf("Dataset Instances: {%d}", len(points))
	log.Print("Timer set")

	// Splitting Dataset
	chunks := splitChunks(points, len(mappers))

	// centroids := startingCentroids(points, in.Clusters) // Standard Centroids Initialization
	centroids := startingCentroidsPlus(points, in.Clusters) // KMeans++ Implementation
	generatedCentroids := centroids
	log.Print("--------------------------------------------------------")
	mChannels, rChannels := initializeChannels()

	defer closeChannels(mChannels, rChannels)

	var reducersReplies []utils.Point
	var iteration int = 1

	for {
		fmt.Print("\n")
		log.Printf("---- Iteration [%d] ----", iteration)

		// Sending Chunk to Mapper
		for index, mapper := range mappers {
			go sendToMapper(chunks[index], centroids, mapper, mChannels[index])
		}
		// sync barrier
		waitMappersResponse(mChannels)

		// Comunicate to reducer which cluster's key has to obtain.
		for index := 0; index < in.Clusters; index++ {
			go sendToReducer(utils.ReducerInput{Mappers: mappers, ClusterKey: index}, reducers[index], rChannels[index])
			if cfg.Parameters.COMBINER {
				// We'll send only to first online reducer if Combiner is ON
				break
			}
		}
		log.Printf("Waiting for replies from reducers\n\n")
		reducersReplies = formalize(waitReducersResponse(rChannels, in.Clusters))
		iteration++

		// logging of centroid calculated at i-iteration
		for i, rep := range reducersReplies {
			log.Printf("New Centroid #%d: %v", i, rep.Values)
		}

		//Checking convergence of Algorithm
		if (checkConvergence(reducersReplies, centroids)) || iteration > 50 {
			centroids = reducersReplies
			break
		}

		centroids = reducersReplies
	}

	reply.Iterations = iteration
	reply.Centroids = centroids
	reply.ExecutionTime = time.Since(start)
	reply.StartingCentroids = generatedCentroids
	log.Print("--------------------------------------------------------")
	log.Print("Convergence achieved, the results were sent to the customer")
	log.Print("--------------------------------------------------------\n")
	return nil

}

// A single Mapper receive a chunk of points and list of actual centroids

func sendToMapper(chunk []utils.Point, centroids []utils.Point, mapper utils.WorkerInfo, ch chan string) {
	addr := mapper.IP + ":" + cfg.Server.WORKER_PORT
	client, err := rpc.Dial("tcp", addr)
	log.Print("Sending chunks to mapper: ", addr)

	if err != nil {
		log.Fatal("Error in dialing with worker: ", err)
	}
	defer client.Close()

	var reply string
	err = client.Call("Mapper.Map", utils.MapperInput{Chunk: chunk, Centroids: centroids}, &reply)
	if err != nil {
		log.Fatal("Error in Mapper.Map: ", err.Error())
	}
	ch <- reply
}

// A single Reducer receive all points mapped to a single cluster

func sendToReducer(input utils.ReducerInput, reducer utils.WorkerInfo, ch chan utils.ReducerResponse) {
	addr := reducer.IP + ":" + cfg.Server.WORKER_PORT
	client, err := rpc.Dial("tcp", addr)
	log.Print("Sending centroid Key to reducer: ", addr)

	if err != nil {
		log.Fatal("Error in dialing with worker: ", err)
	}
	defer client.Close()

	var reply utils.ReducerResponse
	err = client.Call("Reducer.Reduce", input, &reply)
	if err != nil {
		log.Fatal("Error in Reducer.Reduce: ", err.Error())
	}

	ch <- reply

}
