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
	"encoding/csv"
	"fmt"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const DATASET_DIR string = "datasets/"

type Dataset struct {
	Path      string
	Name      string
	Instances int
	K         int
}

// dataset 3d 2cluster 1000samples
func main() {

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
	var fileName string = "d_mapper_" + utils.TEST_FILE
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
	// log.Print("Time Elapsed: ", time.Since(start).Milliseconds())

	if len(os.Args) < 1 {
		fmt.Println("Please specify master address to perform tests")
		os.Exit(1)
	}
	addr := os.Args[1]
	port := ":" + os.Args[2]

	client, err := rpc.DialHTTP("tcp", addr+port)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	// for i, dataset := range datasets {
	for i := 0; i < 4; i++ {
		dataset := datasets[i]
		log.Printf("Test #%d on %s", i, dataset.Name)
		start := time.Now()

		kmeansInput := utils.InputKMeans{Dataset: dataset.Name, Clusters: dataset.K}
		var finalResult utils.Result

		err = client.Call("Master.KMeans", kmeansInput, &finalResult)
		if err != nil {
			log.Fatal("Error in Master.KMeans: \n", err.Error())
		}
		saveResults(finalResult, dataset, time.Since(start), fileName)
	}

}

func saveResults(res utils.Result, dataset Dataset, executionTime time.Duration, fileName string) bool {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer f.Close()

	csvwriter := csv.NewWriter(f)
	defer csvwriter.Flush()

	var line []string
	line = append(line, strconv.Itoa(dataset.Instances), strconv.Itoa(dataset.K), strconv.Itoa(res.Iterations), strconv.Itoa(int(executionTime.Milliseconds())))
	csvwriter.Write(line)
	// }

	// log.Printf("Results are available in [%s]", filePath)
	// csvwriter.Flush()
	return true
}

func TouchFile(name string) os.File {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("File creation error")
	}
	return *file
}
