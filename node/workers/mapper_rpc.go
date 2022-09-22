package main

import (
	"kmeans-MR/utils"
	"math"
	"os"
)

type Mapper int

func (w *Mapper) Map(input utils.MapperInput, reply *utils.MapperResponse) error {

	var dimension int = len(input.Chunk[0].Values)
	var clusters = make([][]utils.Point, 3)
	var minDistance float64 = 0
	var centroidIndex int

	for _, point := range input.Chunk {

		for i := 0; i < len(input.Centroids); i++ {
			euDistance := euclideanDistance(point, input.Centroids[i], dimension) // non serve salvare le distanze, in input ai reducer servono solo i cluster composti
			// log.Printf("Distance from centroid #%d: %f", i, euDistance)

			// first distance calculated should be setted as min (i==0)
			if euDistance <= minDistance || i == 0 {
				// log.Print("La distanza minore è dal centroide #", i)
				minDistance = euDistance
				centroidIndex = i
			}
		}

		clusters[centroidIndex] = append(clusters[centroidIndex], point)
		centroidIndex = 0
		minDistance = 0
	}
	utils.ViewClusters(clusters, len(input.Centroids), false)
	reply.Clusters = clusters
	reply.IP = os.Getenv("HOSTNAME")
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
