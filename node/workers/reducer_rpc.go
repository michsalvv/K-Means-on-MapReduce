package main

import (
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

type Reducer int

var cluster []utils.Point

func (r *Reducer) Reduce(in utils.ReducerInput, reply *utils.ReducerResponse) error {
	log.Printf("Processing Cluster #[%d] ", in.ClusterKey)
	for _, mapper := range in.Mappers {
		clusterPoints := request(mapper, in.ClusterKey).Cluster //TODO se Dial è bloccante e connessioni HTTP non vanno bene, pensare ad una go routine per reducer
		cluster = append(cluster, clusterPoints...)
	}

	log.Printf("Received [%d] points for cluster #%d  ", len(cluster), in.ClusterKey)

	// var newCentroid utils.Point = recenter()
	*reply = utils.ReducerResponse{Centroid: recenter(cluster), IP: os.Getenv("HOSTNAME")}
	return nil
}

func request(mapper utils.WorkerInfo, clusterKey int) utils.MapperResponse {
	addr := mapper.IP + ":" + strconv.Itoa(utils.WORKER_PORT)
	//TODO Vedere se Dial va bene perchè potrebbe essere bloccante quindi un mapper potrebbe rispondere alle richieste dai vari reducer in sequenza e non contemporaneamente
	client, err := rpc.Dial("tcp", addr)
	log.Print("Asking data to mapper: ", addr)

	if err != nil {
		log.Print("Error in dialing with worker: ", err)
	}
	defer client.Close()

	var reply utils.MapperResponse
	err = client.Call("Mapper.GetClusters", clusterKey, &reply)
	if err != nil {
		log.Print("Error in Mapper.Map: ", err.Error())
	}
	return reply
}

func recenter(points []utils.Point) utils.Point {

	var dimension int = len(points[0].Values)
	log.Printf("Recenter: dimension [%d]", dimension)

	centroidValues := make([]float64, dimension)

	for _, point := range points {
		log.Printf("%f", point.Values)
		for i := 0; i < dimension; i++ {
			centroidValues[i] += point.Values[i]
			// log.Print("updated: ", centroidValues)
		}
	}

	log.Print(centroidValues)
	for i := 0; i < dimension; i++ {
		centroidValues[i] = centroidValues[i] / float64(len(points))
	}

	log.Print("Recentered centroid: ", centroidValues)
	utils.Wait()
	return utils.Point{Values: centroidValues}
}
