package main

import (
	"kmeans-MR/utils"
	"log"
	"math"
	"os"
)

type Mapper int

var clustersPointer *[][]utils.Point
var chunk *[]utils.Point
var local_sum []utils.Point
var dimension int

// RPC exposed to perform the Map Task
func (w *Mapper) Map(input utils.MapperInput, reply *string) error {

	var minDistance float64 = 0
	var centroidIndex int
	local_sum = nil

	// Retrieving of chunk ad metadata only at first iteration
	if input.Chunk != nil {
		chunk = &input.Chunk
		dimension = len(input.Chunk[0].Values)
	}

	var clusters = make([][]utils.Point, len(input.Centroids)) // len(input.Centroids) is K
	for _, point := range *chunk {

		for i := 0; i < len(input.Centroids); i++ {
			euDistance := euclideanDistance(point, input.Centroids[i], dimension)
			// first distance calculated should be setted as min (i==0)
			if euDistance <= minDistance || i == 0 {
				minDistance = euDistance
				centroidIndex = i
			}
		}

		clusters[centroidIndex] = append(clusters[centroidIndex], point)
		centroidIndex = 0
		minDistance = 0
	}
	utils.ViewClusters(clusters, len(input.Centroids), false)

	clustersPointer = &clusters
	if cfg.Parameters.COMBINER {
		local_sum = combine(clustersPointer)
	}
	*reply = os.Getenv("HOSTNAME")
	return nil
}

/*
* Compute for each centroid local sums of points
* Sendo to reducer: <centroid, partial sums>
* Used only if Combiner=ON
 */
func combine(clusters *[][]utils.Point) []utils.Point {
	var combined_values = make([]utils.Point, len(*clusters))

	for j, cluster := range *clusters {
		centroidValues := make([]float64, dimension)

		for _, point := range cluster {
			for i := 0; i < dimension; i++ {
				centroidValues[i] += point.Values[i]
			}
		}
		combined_values[j].Values = centroidValues
	}
	return combined_values
}

// RPC exposed to Reducers to retrieve data
func (w *Mapper) GetClusters(input int, reply *utils.MapperResponse) error {
	log.Print("Request recieved from reducer with clusterKey: ", input)

	if cfg.Parameters.COMBINER {
		clusterDimensionality := make([]int, len(*clustersPointer)) //len(*clusterPointer) -> k
		for i := 0; i < len(*clustersPointer); i++ {
			clusterDimensionality[i] = len((*clustersPointer)[i])
		}
		*reply = utils.MapperResponse{IP: os.Getenv("HOSTNAME"), Cluster: local_sum, ClusterDimensionality: clusterDimensionality}
		return nil
	}
	*reply = utils.MapperResponse{Cluster: (*clustersPointer)[input], IP: os.Getenv("HOSTNAME")}
	return nil
}

// Square root of the sums of the swuare of the differences between the coordinates of the points in each dimension
func euclideanDistance(point, centroid utils.Point, d int) float64 {
	var distance float64
	pointVals := point.Values
	centroidVals := centroid.Values

	for i := 0; i < d; i++ {
		distance += math.Pow(pointVals[i]-centroidVals[i], 2)
	}
	return math.Sqrt(distance)
}
