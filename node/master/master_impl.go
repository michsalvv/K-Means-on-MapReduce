package main

import (
	"bufio"
	"kmeans-MR/utils"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func startingCentroids(points []utils.Point, kValue int) []utils.Point {
	centroids := make([]utils.Point, kValue)
	rand.Seed(time.Now().UnixNano()) // Initialization of the source used from rand

	for i := 0; i < kValue; i++ {
		randIndex := rand.Intn(len(points))
		log.Print("randIndex: ", randIndex)
		centroids[i] = points[randIndex]
	}
	log.Print("Starting Centroids: ", centroids)
	return centroids
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

func sortResponse(channels map[int]chan [][]utils.Point) [][]utils.Point {
	var s [][]utils.Point
	// Waiting for #Workers replies
	for i := 0; i < len(channels); i++ {
		s = <-channels[i]
	}
	return s
}
