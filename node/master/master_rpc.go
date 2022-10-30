package main

import (
	"bufio"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

// const masterPath string = "server/master/"

var mappers []utils.WorkerInfo
var reducers []utils.WorkerInfo

type Master int

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

func (m *Master) KMeans(in utils.InputKMeans, reply *utils.Result) error {

	clusterError := checkAvailability(in, mappers, reducers)

	if clusterError != nil {
		return clusterError
	}

	log.Print("=========== START KMEANS ===========")
	log.Printf("K-Means on {%s}: {%d} clusters", in.Dataset, in.Clusters)

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

	chunks := splitChunks(points, len(mappers))
	// centroids := startingCentroids(points, in.Clusters)

	centroids := startingCentroidsPlus(points, in.Clusters)

	mChannels, rChannels := initializeChannels()

	defer closeChannels(mChannels, rChannels)

	// First iteration
	for index, mapper := range mappers {
		go sendToMapper(chunks[index], centroids, mapper, mChannels[index])
	}

	// TODO aggiungi controllo sul bool di waitMappersResponse
	waitMappersResponse(mChannels) // Funge da barriera di sincronizzazione

	var prevCentroids []utils.Point
	var reducersReplies []utils.Point
	var iteration int = 1

	//TODO mettere la prima iterazione dentro il ciclo
	//TODO inverti ordine ciclo
	for {

		log.Printf("---- ITERATION [%d] ----", iteration)
		// Comunicate to reducer which cluster's key has to obtain.

		for index := 0; index < in.Clusters; index++ {
			go sendToReducer(utils.ReducerInput{Mappers: mappers, ClusterKey: index}, reducers[index], rChannels[index])
		}
		reducersReplies = formalize(waitReducersResponse(rChannels, in.Clusters)) //trova nome migliore di formalize

		// logging of centroid calculated at i-iteration
		for i, rep := range reducersReplies {
			log.Print("Centroid #", i, ": ", rep)
		}
		// At the beginning prevCentroids = nil
		if len(prevCentroids) != 0 {
			if convergence(reducersReplies, prevCentroids) {
				break
			}
		}
		prevCentroids = reducersReplies

		for index, mapper := range mappers {
			go sendToMapper(nil, reducersReplies, mapper, mChannels[index])
		}

		waitMappersResponse(mChannels)

		iteration++
	}

	reply.Iterations = iteration
	reply.Centroids = reducersReplies

	return nil

}

// A single Mapper receive a chunk of points and list of actual centroids

func sendToMapper(chunk []utils.Point, centroids []utils.Point, mapper utils.WorkerInfo, ch chan string) {
	addr := mapper.IP + ":" + strconv.Itoa(utils.WORKER_PORT)
	client, err := rpc.Dial("tcp", addr)
	log.Print("Sending data to mapper: ", addr)

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
	addr := reducer.IP + ":" + strconv.Itoa(utils.WORKER_PORT)
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
