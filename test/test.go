package main

import (
	"fmt"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"

	"github.com/fatih/color"
)

var client *rpc.Client
var err error

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify service [address:port] [MULTIPLE/ALLDATA] [numerOfMappers] optional{[datasetName], [RUNS NUMBER] to perform tests")
		os.Exit(1)
	}

	// Connecting to service
	addr := os.Args[1]
	MAPPER_NUMS = os.Args[3]
	client, err = rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	switch os.Args[2] {
	case MODE_ALL_DATASET:
		run_test(os.Args[4], MODE_ALL_DATASET, 1)

	case MODE_MULTIPLE:
		var runs int
		if len(os.Args) > 4 {
			if runs, _ = strconv.Atoi(os.Args[5]); runs > 0 {
				PrintResults(run_test(os.Args[4], MODE_MULTIPLE, runs))
			}
			break
		}
		log.Fatal("Please use a valid run test numbers!")

	default:
		log.Fatal("Please select a valid test mode [MULTIPLE/ALLDATA]")
	}
}

func run_test(datasetName, testingMode string, runs int) []utils.Result {

	var results []utils.Result

	switch testingMode {

	case MODE_MULTIPLE:
		dataset := FetchSingleDataset(datasetName)
		file := TouchFile(MAPPER_NUMS + TEST_FILE_FORMAT)
		log.Printf("Testing on %s", dataset.Name)
		for i := 0; i < runs; i++ {
			kmeansInput := utils.InputKMeans{Dataset: dataset.Name, Clusters: dataset.K}
			results = append(results, KMeans(kmeansInput))
		}
		SaveBenchmark(results, dataset, file.Name())

	case MODE_ALL_DATASET:
		file := TouchFile("ALL_DATA_" + MAPPER_NUMS + TEST_FILE_FORMAT)
		datasets := fetchDatasets()

		for _, dataset := range datasets {
			kmeansInput := utils.InputKMeans{Dataset: dataset.Name, Clusters: dataset.K}
			finalResult := KMeans(kmeansInput)

			results = append(results, finalResult)
			SaveBenchmark([]utils.Result{finalResult}, dataset, file.Name())
			log.Printf("[%s] results achieved in [%s] iterations with [%s] using [%s] mapper nodes",
				color.HiWhiteString(dataset.Name),
				color.YellowString(strconv.Itoa(finalResult.Iterations)),
				color.GreenString(finalResult.ExecutionTime.String()), color.RedString(MAPPER_NUMS))
		}
	}

	return results
}
