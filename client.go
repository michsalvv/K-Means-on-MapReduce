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

	"github.com/fatih/color"
)

const CLIENT_OUT_DIR string = "datasets/"

func main() {
	cfg := utils.GetConfiguration()
	if len(os.Args) < 2 {
		fmt.Println("Please specify datasetPath and number of cluster to find:\n\tgo run client.go [dataset] [#clusters]")
		os.Exit(1)
	}
	addr := cfg.Server.HOST + ":" + cfg.Server.MASTER_PORT
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	clusters, _ := strconv.Atoi(os.Args[2])
	datasetPath := strings.Split(os.Args[1], "/")
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
	log.Print(color.HiWhiteString("------- Clustering via KMeans of "), color.HiGreenString("{"+datasetName+"}"), color.HiWhiteString(" -------"))
	log.Print("\n")
	log.Print("Starting Centroids:")
	for _, centroid := range res.StartingCentroids {
		line := strings.Fields(strings.Trim(fmt.Sprint(centroid), "{}"))
		log.Printf("%s", color.HiBlueString(strings.Join(line, "; ")))
	}

	log.Printf("Convergence achieved in [%s] %s: ", color.HiGreenString(strconv.Itoa(res.Iterations)), color.HiGreenString("iterations"))
	for i, point := range res.Centroids {
		line := strings.Fields(strings.Trim(fmt.Sprint(point.Values), "[]"))
		log.Printf("[#%s] ->  %s", color.HiYellowString(strconv.Itoa(i)), color.HiYellowString(strings.Join(line, "; ")))
		csvwriter.Write(line)
	}

	fmt.Println()
	log.Printf("Results are available in [%s]", color.HiRedString(filePath))
	csvwriter.Flush()
	return true
}
