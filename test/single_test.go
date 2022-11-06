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
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify service [address:port] and [numerOfMappers] to perform tests")
		os.Exit(1)
	}
	addr := os.Args[1]

	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	var datasets []Dataset
	datasets, fileName := fetchDatasets()

	for i, dataset := range datasets {
		log.Printf("Test #%d on %s", i, dataset.Name)
		kmeansInput := utils.InputKMeans{Dataset: dataset.Name, Clusters: dataset.K}
		var finalResult utils.Result

		err = client.Call("Master.KMeans", kmeansInput, &finalResult)
		if err != nil {
			log.Fatal("Error in Master.KMeans: \n", err.Error())
		}
		if saveBenchmark(finalResult, dataset, fileName) {
			log.Print("Test Done ...")
			saveResults(finalResult, dataset.Name)
		}
	}
}

// dataset 3d 2cluster 1000samples format
func fetchDataset() ([]Dataset, string) {
	log.Print("Fetching Datasets ...")
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	exPath := filepath.Dir(ex)
	datasetsPath := strings.Replace(exPath, "test", DATASET_DIR, 1)

	files, err := os.ReadDir(datasetsPath)
	if err != nil {
		log.Fatal(err)
	}
	var fileName string = os.Args[2] + "_mapper_" + utils.TEST_FILE
	TouchFile(fileName)

	var datasets []Dataset
	for _, dir := range files {
		tokens := strings.Split(dir.Name(), "_")
		re := regexp.MustCompile("[0-9]+")
		if dir.IsDir() && strings.Contains(tokens[0], "dataset") {
			k, _ := strconv.Atoi(re.FindString(tokens[2]))
			inst, _ := strconv.Atoi(re.FindString(tokens[3]))
			datasets = append(datasets, Dataset{Path: DATASET_DIR + dir.Name(), K: k, Instances: inst, Name: dir.Name() + ".csv"})
		}
	}

	return datasets, fileName
}
