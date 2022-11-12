package main

import (
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
)

type Reducer int

func (r *Reducer) Reduce(in utils.ReducerInput, reply *utils.ReducerResponse) error {
	log.Printf("Processing Cluster #[%d] ", in.ClusterKey)

	var cluster []utils.Point
	for _, mapper := range in.Mappers {
		clusterPoints := request(mapper, in.ClusterKey).Cluster
		cluster = append(cluster, clusterPoints...)
	}

	log.Printf("Received [%d] points for cluster #%d  ", len(cluster), in.ClusterKey)

	*reply = utils.ReducerResponse{Centroid: recenter(cluster), IP: os.Getenv("HOSTNAME")}
	return nil
}

func request(mapper utils.WorkerInfo, clusterKey int) utils.MapperResponse {
	addr := mapper.IP + ":" + cfg.Server.WORKER_PORT
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
	centroidValues := make([]float64, dimension)

	for _, point := range points {
		for i := 0; i < dimension; i++ {
			centroidValues[i] += point.Values[i]
		}
	}
	for i := 0; i < dimension; i++ {
		centroidValues[i] = centroidValues[i] / float64(len(points))
	}

	log.Print("Recentered centroid: ", centroidValues)
	return utils.Point{Values: centroidValues}
}
