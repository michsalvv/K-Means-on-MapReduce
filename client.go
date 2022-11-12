package main

import (
	"encoding/csv"
	"fmt"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

const CLIENT_OUT_DIR string = "datasets/"

func main() {
	cfg := utils.GetConfiguration()
	fmt.Printf("%+v", cfg)
	if len(os.Args) < 3 {
		fmt.Println("Please specify master address, dataset and number of cluster to find:\n\tgo run client.go [master] [dataset] [#clusters]")
		os.Exit(1)
	}
	// addr := os.Args[1]
	addr := cfg.Server.HOST + ":" + cfg.Server.MASTER_PORT
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	clusters, _ := strconv.Atoi(os.Args[3])
	datasetPath := strings.Split(os.Args[2], "/")
	datasetName := datasetPath[len(datasetPath)-1]

	kmeansInput := utils.InputKMeans{Dataset: datasetName, Clusters: clusters}
	var finalResult utils.Result

	err = client.Call("Master.KMeans", kmeansInput, &finalResult)
	if err != nil {
		log.Fatal("Error in Master.KMeans: \n", err.Error())
	}
	if saveResults(finalResult, datasetName) {
		log.Printf("Use client/check_results.py to validate the results")
	}
}

func saveResults(res utils.Result, datasetName string) bool {
	filename := strings.Replace(datasetName, "dataset", "centroids", 1)
	filePath := fmt.Sprintf("%s%s%s", CLIENT_OUT_DIR, strings.Replace(datasetName, ".csv", "/", 1), filename)

	csvFile, err := os.Create(filePath)
	if err != nil {
		log.Print("Failed creating file: ", err)
		return false
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	fmt.Println()
	log.Printf("Convergence achieved in %d iterations: ", res.Iterations)
	for i, point := range res.Centroids {
		line := strings.Fields(strings.Trim(fmt.Sprint(point.Values), "[]"))
		log.Printf("[#%d] ->  %s", i, strings.Join(line, "; "))
		csvwriter.Write(line)
	}

	fmt.Println()
	log.Printf("Results are available in [%s]", filePath)
	csvwriter.Flush()
	return true
}
