package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

const MASTER_IP string = "localhost"

type Input struct {
	Text       string
	WordToGrep string
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify file and word:\n\tgo run client.go [word] [file]")
		os.Exit(1)
	}
	addr := MASTER_IP + ":" + strconv.Itoa(9001)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	grepInput := Input{Text: os.Args[2], WordToGrep: os.Args[1]}
	var greppedText string
	err = client.Call("Master.Grep", grepInput, &greppedText)
	if err != nil {
		log.Fatal("Error in Master.Grep: ", err.Error())
	}

	log.Println("MASTER GREP REPLY\n", greppedText)
}
