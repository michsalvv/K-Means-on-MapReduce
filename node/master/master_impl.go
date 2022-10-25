package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"kmeans-MR/utils"
	"log"
	"math"
	"math/big"

	"strconv"
	"strings"
)

func startingCentroids(points []utils.Point, kValue int) []utils.Point {
	centroids := make([]utils.Point, kValue)

	for i := 0; i < kValue; i++ {
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(points))))
		log.Print("randIndex: ", randIndex)
		centroids[i] = points[randIndex.Int64()]
	}
	log.Print("Starting Centroids: ", centroids)
	return centroids
}

// TODO rimuovere il centroide dall'insieme dei punti una volta scelto, perchè se viene scelto lo stesso centroide due volte, ad un reducer non verrà inviato nulla
func startingCentroidsPlus(points []utils.Point, kValue int) []utils.Point {
	dimension := len(points[0].Values)
	centroids := make([]utils.Point, kValue)

	randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(points))))

	// The first centroid is selected randomly
	var foundedCentroid int = 0
	centroids[foundedCentroid] = points[randIndex.Int64()]
	foundedCentroid++

	var maxDist float64
	var farthestPoint utils.Point

	// 2nd step: the point farthest away will become next centroid
	for _, point := range points {
		dist := euclideanDistance(point, centroids[0], dimension)
		if dist >= maxDist {
			maxDist = dist
			farthestPoint = point
		}
	}
	log.Print("dist: ", maxDist)
	log.Print("farthestPoint: ", farthestPoint)

	centroids[foundedCentroid] = farthestPoint
	foundedCentroid++

	// Iterate untile k centroids has been selected
	var iteration int = 0
	for {
		log.Print(iteration)
		iteration++
		var clusters = make([][]utils.Point, kValue)
		maxDist = 0
		for _, point := range points {
			var minDistance float64 = 0
			var centroidIndex int = 0

			for i := 0; i < foundedCentroid; i++ {
				euDistance := euclideanDistance(point, centroids[i], dimension)

				if euDistance <= minDistance || i == 0 {
					minDistance = euDistance
					centroidIndex = i
				}
				if euDistance >= maxDist {
					maxDist = euDistance
					farthestPoint = point
				}
			}

			clusters[centroidIndex] = append(clusters[centroidIndex], point)
			centroidIndex = 0
			minDistance = 0
		}

		centroids[foundedCentroid] = farthestPoint
		foundedCentroid++ // must be updated before check beacuse it's start from 0
		if foundedCentroid == kValue {
			break
		}
	}
	log.Print(centroids)
	return centroids
}

func euclideanDistance(point, centroid utils.Point, d int) float64 {
	var distance float64
	pointVals := point.Values
	centroidVals := centroid.Values

	for i := 0; i < d; i++ {
		distance += math.Pow(pointVals[i]-centroidVals[i], 2)
	}
	return math.Sqrt(distance)
}

func formalize(replies []utils.ReducerResponse) []utils.Point {
	var ret []utils.Point
	for _, rep := range replies {
		ret = append(ret, rep.Centroid)
	}

	return ret
}

func convergence(actual, prev []utils.Point) bool {
	var ratio float64
	var dimension int = len(actual[0].Values)

	for i, point := range actual {
		for j := 0; j < dimension; j++ {
			ratio = point.Values[j] - prev[i].Values[j]
			log.Print("CONVERGENCE: ", ratio)
			if ratio > 0.001 { //TODO define configurable threshold
				return false
			}
		}
	}

	return true

}

func splitChunks(points []utils.Point, numberChunks int) [][]utils.Point {
	x := int(float64(len(points)) / float64(numberChunks))
	log.Printf("[%d Points]\t[%d Mappers]\t%d points for mapper", len(points), numberChunks, x)

	var offset, endOffset int
	splittedChunks := make([][]utils.Point, numberChunks)
	for i := 0; i < numberChunks; i++ {
		offset = i * x
		if (i == numberChunks-1) && (numberChunks > 1) {
			endOffset = len(points)
		} else {
			endOffset = offset + x
		}
		// log.Printf("Mapper [%d]\tfrom: %d \tto: %d\n", i, offset, endOffset)
		for j := offset; j < endOffset; j++ {
			splittedChunks[i] = append(splittedChunks[i], points[j])
		}
	}
	return splittedChunks
}

func readPoint(r *bufio.Reader) (utils.Point, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line     []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
	}

	if len(line) == 0 {
		return utils.Point{Values: nil}, err
	}

	var values []float64
	for _, val := range strings.Split(string(line), ",") {
		floated, _ := strconv.ParseFloat(val, 64)
		values = append(values, floated)
	}
	return utils.Point{Values: values}, err
}

// Barriera di sincronizzazione
func waitMappersResponse(channels map[int]chan string) bool {
	var replies []string

	// Waiting for #Workers replies
	for i := 0; i < len(channels); i++ {
		replies = append(replies, <-channels[i])
	}

	log.Print("All the mappers responded")

	return true
}

func waitReducersResponse(channels map[int]chan utils.ReducerResponse, dutyReducers int) []utils.ReducerResponse {
	var replies []utils.ReducerResponse

	// Waiting for reducers replies
	for i := 0; i < dutyReducers; i++ {
		replies = append(replies, <-channels[i])
	}

	log.Print("All the reducers responded")

	return replies
}

func initializeChannels() (map[int]chan string, map[int]chan utils.ReducerResponse) {

	mChannels := make(map[int]chan string)
	rChannels := make(map[int]chan utils.ReducerResponse)

	for index := range mappers {
		mChannels[index] = make(chan string)
	}

	for index := range reducers {
		rChannels[index] = make(chan utils.ReducerResponse)
	}

	return mChannels, rChannels
}

func closeChannels(mChannels map[int]chan string, rChannels map[int]chan utils.ReducerResponse) {
	for index := range mappers {
		close(mChannels[index])
	}

	for index := range reducers {
		rChannels[index] = make(chan utils.ReducerResponse)
		close(rChannels[index])
	}
}

func checkAvailability(inputData utils.InputKMeans, mappers, reducers []utils.WorkerInfo) error {
	if len(mappers) == 0 || len(reducers) == 0 {
		log.Print(utils.NO_RES_ERROR)
		return errors.New(utils.NO_RES_ERROR)
	}

	if len(reducers) < inputData.Clusters {
		log.Print(utils.NO_REDUCERS_ERROR)
		return errors.New(utils.NO_REDUCERS_ERROR)
	}

	return nil
}
