package main

import (
	"encoding/csv"
	"fmt"
	"kmeans-MR/utils"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const TEST_FILE_FORMAT = "_mapper_results.csv"

type Dataset struct {
	Path      string
	Name      string
	Instances int
	K         int
}

const DATASET_DIR string = "datasets"
const MODE_ALL_DATASET = "ALLDATA"
const MODE_MULTIPLE = "MULTIPLE"

var MAPPER_NUMS string

func KMeans(kmeansInput utils.InputKMeans) utils.Result {
	var finalResult utils.Result

	err := client.Call("Master.KMeans", kmeansInput, &finalResult)
	if err != nil {
		log.Fatal("Error in Master.KMeans: \n", err.Error())
	}
	return finalResult
}

func fetchDatasets() []Dataset {
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
	return datasets
}

/*
Dataset name format {dataset_dimension_clusters_instances.csv}
*/
func FetchSingleDataset(datasetName string) Dataset {
	datasetName = strings.Split(datasetName, ".")[0]
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	exPath := filepath.Dir(ex)
	datasetPath := fmt.Sprintf("%s/%s/%s.csv", strings.Replace(exPath, "test", DATASET_DIR, 1), datasetName, datasetName)

	/*
		check if file exists and takes metadata
	*/
	var dataset Dataset
	re := regexp.MustCompile("[0-9]+")
	if _, err := os.Stat(datasetPath); err == nil {
		tokens := strings.Split(datasetName, "_")
		k, _ := strconv.Atoi(re.FindString(tokens[2]))
		inst, _ := strconv.Atoi(re.FindString(tokens[3]))
		dataset = Dataset{Path: datasetPath, K: k, Instances: inst, Name: datasetName + ".csv"}
	} else {
		log.Fatal("Please use a valid dataset name. Include .csv extension!")
	}

	return dataset
}

func TouchFile(name string) os.File {
	exists, _ := os.Stat(name)
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("File creation error")
	}
	defer file.Close()
	csvwriter := csv.NewWriter(file)

	if exists == nil {
		var header []string
		header = append(header, "Dataset", "Points", "Clusters", "Iterations", "ExecutionTime[ms]", "Mappers")
		csvwriter.Write(header)
		csvwriter.Flush()
	}
	return *file
}

func SaveBenchmark(res []utils.Result, dataset Dataset, fileName string) bool {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer f.Close()

	csvwriter := csv.NewWriter(f)
	defer csvwriter.Flush()
	for _, result := range res {
		var line []string
		executionTime := result.ExecutionTime.Milliseconds()
		line = append(line, dataset.Name, strconv.Itoa(dataset.Instances), strconv.Itoa(dataset.K),
			strconv.Itoa(result.Iterations), strconv.FormatInt(executionTime, 10), MAPPER_NUMS)
		csvwriter.Write(line)
	}
	return true
}

func SaveResults(res utils.Result, datasetName string) bool {
	filename := strings.Replace(datasetName, "dataset", "centroids", 1)

	csvFile, err := os.Create(filename)
	if err != nil {
		log.Print("Failed creating file: ", err)
		return false
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	for _, point := range res.Centroids {
		line := strings.Fields(strings.Trim(fmt.Sprint(point.Values), "[]"))
		csvwriter.Write(line)
	}
	csvwriter.Flush()
	return true
}

func PrintResults(results []utils.Result) {
	var mean time.Duration

	log.Print(color.HiGreenString("Test done!"))

	for i, res := range results {
		var executionTime string = res.ExecutionTime.String()
		var iterations string = strconv.Itoa(res.Iterations)
		log.Printf("[%s] results achieved in [%s] iterations with [%s] using [%s] mapper nodes",
			color.HiWhiteString("Iteration #"+strconv.Itoa(i)),
			color.YellowString(iterations), color.GreenString(executionTime), color.RedString(MAPPER_NUMS))
		mean += res.ExecutionTime
	}

	/*
		Print calculated centroids of only first run
	*/
	log.Print("Results:")
	for i, point := range results[0].Centroids {
		line := strings.Fields(strings.Trim(fmt.Sprint(point.Values), "[]"))
		log.Printf("RUN [#%d] ->  %s", i, strings.Join(line, "; "))
	}

	if len(results) > 1 {
		/*
			len(results) = test runs
		*/
		mean = mean / time.Duration(len(results))
		log.Print(color.HiCyanString("AVERAGE EXECUTION TIME: [%s]", mean))
	}
}
