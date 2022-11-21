package main

import (
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
)

type Reducer int

var pointsOfCluster []int

func (r *Reducer) Reduce(in utils.ReducerInput, reply *utils.ReducerResponse) error {
	log.Printf("Processing Cluster #[%d] ", in.ClusterKey)

	var cluster []utils.Point
	pointsOfCluster = nil
	for i, mapper := range in.Mappers {

		if cfg.Parameters.COMBINER {
			combined_response := retrieveData(mapper, 0)

			if i == 0 {
				cluster = make([]utils.Point, len(combined_response.Cluster))
				for i := 0; i < len(combined_response.Cluster); i++ {
					cluster[i].Values = make([]float64, len(combined_response.Cluster[0].Values))
				}
			}
			aggregate(&cluster, combined_response.Cluster, combined_response.ClusterDimensionality)

		} else {
			clusterPoints := retrieveData(mapper, in.ClusterKey).Cluster
			cluster = append(cluster, clusterPoints...)
		}
	}

	log.Printf("Received [%d] points for cluster #%d  ", len(cluster), in.ClusterKey)
	recenteredCluster := recenter(cluster)
	if cfg.Parameters.COMBINER {
		*reply = utils.ReducerResponse{CombinedResponse: recenteredCluster, IP: os.Getenv("HOSTNAME")}
		return nil
	}
	*reply = utils.ReducerResponse{Centroid: recenteredCluster[0], IP: os.Getenv("HOSTNAME")}
	return nil
}

func retrieveData(mapper utils.WorkerInfo, clusterKey int) utils.MapperResponse {
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

func recenter(points []utils.Point) []utils.Point {

	var dimension int = len(points[0].Values)
	clusters := points
	if cfg.Parameters.COMBINER {
		for k, cluster := range clusters {
			for i := 0; i < dimension; i++ {
				cluster.Values[i] = cluster.Values[i] / float64(pointsOfCluster[k])
			}
		}
		return clusters
	}
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
	result := utils.Point{Values: centroidValues}
	return []utils.Point{result}
}

func aggregate(combined *[]utils.Point, localsum []utils.Point, clusterDimensionality []int) {
	dim := len((*combined)[0].Values)

	if pointsOfCluster == nil {
		pointsOfCluster = make([]int, len(clusterDimensionality))
	}

	for k, cluster := range *combined {
		for i := 0; i < dim; i++ {
			cluster.Values[i] += localsum[k].Values[i]
		}
		pointsOfCluster[k] += clusterDimensionality[k]
	}
}
