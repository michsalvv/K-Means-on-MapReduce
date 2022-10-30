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
var dimension int

func (w *Mapper) Map(input utils.MapperInput, reply *string) error {

	var minDistance float64 = 0
	var centroidIndex int

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
	*reply = os.Getenv("HOSTNAME")
	return nil
}

func (w *Mapper) GetClusters(input int, reply *utils.MapperResponse) error {
	log.Print("Request recieved from reducer with clusterKey: ", input)
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
