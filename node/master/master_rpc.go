package main

import (
	"bufio"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"
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

func (m *Master) KMeans(in utils.InputKMeans, reply *string) error {
	if len(mappers) == 0 || len(reducers) == 0 {
		log.Print("Not possible to perform K-Means: insufficient workers online")
		return nil
	}
	log.Print("=========== START KMEANS ===========")
	log.Printf("K-Means on {%s}: {%d} clusters", in.Dataset, in.Clusters)

	file, err := os.Open(in.Dataset)
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
	centroids := startingCentroids(points, in.Clusters)

	channels := make(map[int]chan [][]utils.Point)
	for index, mapper := range mappers {
		channels[index] = make(chan [][]utils.Point)
		defer close(channels[index])
		go sendToMapper(chunks[index], centroids, mapper, channels[index])
	}
	sortResponse(channels) // Funge da barriera di sincronizzazione (forse)

	// ch := make(chan string)
	// defer close(ch)
	// go sendToReduce(response, reducers[0], ch)
	// *reply = <-ch
	return nil
}

func sendToMapper(chunk []utils.Point, centroids []utils.Point, mapper utils.WorkerInfo, ch chan [][]utils.Point) {
	addr := mapper.IP + ":" + strconv.Itoa(utils.WORKER_PORT)
	client, err := rpc.Dial("tcp", addr)
	log.Print("Sendding data to mapper: ", addr)

	if err != nil {
		log.Fatal("Error in dialing with worker: ", err)
	}
	defer client.Close()

	var reply [][]utils.Point
	err = client.Call("Mapper.Map", utils.MapperInput{Chunk: chunk, Centroids: centroids}, &reply)
	if err != nil {
		log.Fatal("Error in Mapper.Map: ", err.Error())
	}

	ch <- reply
}

// func sendToReduce(s string, reducer utils.WorkerInfo, ch chan string) {
// 	addr := reducer.IP + ":" + strconv.Itoa(utils.WORKER_PORT)
// 	client, err := rpc.Dial("tcp", addr)
// 	if err != nil {
// 		log.Fatal("Error in dialing with worker: ", err)
// 	}
// 	defer client.Close()

// 	var reply string
// 	err = client.Call("Reducer.Reduce", s, &reply)
// 	if err != nil {
// 		log.Fatal("Error in Reducer.Reduce: ", err.Error())
// 	}

// 	ch <- reply

// }
