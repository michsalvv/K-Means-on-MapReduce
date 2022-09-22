package main

import (
	"fmt"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify master address, dataset and number of cluster to find:\n\tgo run client.go [master] [dataset] [#clusters]")
		os.Exit(1)
	}
	addr := os.Args[1]
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	clusters, _ := strconv.Atoi(os.Args[3])

	kmeansInput := utils.InputKMeans{Dataset: os.Args[2], Clusters: clusters}

	var greppedText string
	err = client.Call("Master.KMeans", kmeansInput, &greppedText)
	if err != nil {
		log.Fatal("Error in Master.Grep: ", err.Error())
	}

	log.Println("MASTER GREP REPLY\n", greppedText)
}
