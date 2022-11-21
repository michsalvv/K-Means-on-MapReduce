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
		centroids[i] = points[randIndex.Int64()]
	}
	log.Print("Standard Centroids Initialization: ", centroids)
	return centroids
}

func startingCentroidsPlus(points []utils.Point, kValue int) []utils.Point {
	dimension := len(points[0].Values)
	centroids := make([]utils.Point, kValue)

	randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(points))))

	// The first centroid is selected randomly
	var foundedCentroid int = 0
	var index int = int(randIndex.Int64())
	centroids[foundedCentroid] = points[index]
	foundedCentroid++
	points = remove(points, int(index))

	var maxDist float64
	var farthestPoint utils.Point

	// 2nd step: the point farthest away will become next centroid
	for i, point := range points {
		dist := euclideanDistance(point, centroids[0], dimension)
		if dist >= maxDist {
			maxDist = dist
			farthestPoint = point
			index = i
		}
	}

	centroids[foundedCentroid] = farthestPoint
	foundedCentroid++
	points = remove(points, index)

	// Iterate untile k centroids has been selected
	var iteration int = 0
	for {
		if foundedCentroid == kValue {
			break
		}
		iteration++
		var distances = make([]utils.Triple, len(points))
		maxDist = 0

		// Associate every point the distance from his nearest centroid
		for j, point := range points {
			var minDistance float64 = 0
			var centroidIndex int = 0
			for i := 0; i < foundedCentroid; i++ {
				euDistance := euclideanDistance(point, centroids[i], dimension)
				if euDistance <= minDistance || i == 0 {
					minDistance = euDistance
					centroidIndex = i
				}
			}
			distances[j] = utils.Triple{P: point, Distance: minDistance, Centroid: centroids[centroidIndex]}
			centroidIndex = 0
			minDistance = 0
		}

		/*
		 * Select the next centroid from the data points such that the probability of choosing a point as centroid is directly
		 * proportional to its distance from the nearest, previously chosen centroid.
		 */
		for _, info := range distances {
			if info.Distance > maxDist {
				maxDist = info.Distance
				farthestPoint = info.P
			}
		}
		centroids[foundedCentroid] = farthestPoint
		foundedCentroid++
	}
	log.Print("KMeans++ Centroids Initialization: ", centroids)
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

func remove(slice []utils.Point, position int) []utils.Point {
	return append(slice[:position], slice[position+1:]...)
}

func formalize(replies []utils.ReducerResponse) []utils.Point {

	if cfg.Parameters.COMBINER {
		return replies[0].CombinedResponse
	}
	var ret []utils.Point
	for _, rep := range replies {
		ret = append(ret, rep.Centroid)
	}

	return ret
}

func checkConvergence(actual, prev []utils.Point) bool {
	var diff float64
	var dimension int = len(actual[0].Values)
	for i, point := range actual {
		for j := 0; j < dimension; j++ {
			diff = point.Values[j] - prev[i].Values[j]
			if diff > cfg.Parameters.CONV_THRESH {
				log.Print(diff)
				return false
			}
		}
	}
	return true
}

func splitChunks(points []utils.Point, numberChunks int) [][]utils.Point {
	x := int(float64(len(points)) / float64(numberChunks))
	log.Printf("Splitting {%d} points into {%d} chunks...", len(points), numberChunks)

	var offset, endOffset int
	splittedChunks := make([][]utils.Point, numberChunks)
	for i := 0; i < numberChunks; i++ {
		offset = i * x
		if (i == numberChunks-1) && (numberChunks > 1) {
			endOffset = len(points)
		} else {
			endOffset = offset + x
		}
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
	log.Print("All mappers responded\n\n")
	return true
}

func waitReducersResponse(channels map[int]chan utils.ReducerResponse, dutyReducers int) []utils.ReducerResponse {
	var replies []utils.ReducerResponse

	// Waiting for reducers replies
	for i := 0; i < dutyReducers; i++ {
		replies = append(replies, <-channels[i])
		if cfg.Parameters.COMBINER {
			break // we had to wait for only one reducer
		}
	}
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
		return errors.New(utils.NO_RES_ERROR)
	}

	if (len(reducers) < inputData.Clusters) && !cfg.Parameters.COMBINER {
		return errors.New(utils.NO_REDUCERS_ERROR)
	}

	return nil
}
