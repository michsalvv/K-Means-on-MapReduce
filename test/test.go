package main

import (
	/*"encoding/csv"
	"fmt"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"*/

	"fmt"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"
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

	// var datasets []Dataset
	// datasets, fileName := fetchDatasets()

	// for i, dataset := range datasets {
	// 	log.Printf("Test #%d on %s", i, dataset.Name)
	// 	kmeansInput := utils.InputKMeans{Dataset: dataset.Name, Clusters: dataset.K}
	// 	var finalResult utils.Result

	// 	err = client.Call("Master.KMeans", kmeansInput, &finalResult)
	// 	if err != nil {
	// 		log.Fatal("Error in Master.KMeans: \n", err.Error())
	// 	}
	// 	if saveBenchmark(finalResult, dataset, fileName) {
	// 		log.Print("Test Done ...")
	// 		saveResults(finalResult, dataset.Name)
	// 	}
	// }
}

func run_test(datasetName, testingMode string, runs int) []utils.Result {

	dataset := FetchSingleDataset(datasetName)
	var results []utils.Result

	switch testingMode {

	case MODE_MULTIPLE:
		log.Printf("Testing on %s", dataset.Name)
		for i := 0; i < runs; i++ {
			kmeansInput := utils.InputKMeans{Dataset: dataset.Name, Clusters: dataset.K}
			results = append(results, KMeans(kmeansInput))
		}
	}

	return results
}
