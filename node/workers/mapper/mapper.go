package main

import (
	"fmt"
	"kmeans-MR/utils"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify master address and mapper port:\n\t./mapper [addr:port] [port]")
		os.Exit(1)
	}
	addr := os.Args[1]
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	port, _ := strconv.Atoi(os.Args[2])

	var id int
	err = client.Call("Master.JoinGrep", utils.JoinRequest{Port: port, Type: utils.MapperType}, &id)
	if err != nil {
		log.Fatal("Error in Master.JoinGrep: ", err.Error())
	}

	if id != 1 {
		log.Fatal("Request declined")
		os.Exit(1)
	} else {
		log.Println("Request accepted")
	}

	worker := new(Worker)
	server := rpc.NewServer()

	err = server.Register(worker)
	if err != nil {
		log.Fatal("Error on Register(worker): ", err)
	}

	var address string = ":" + strconv.Itoa(port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error in listening:", err)
	}
	log.Printf("Worker online on port %d\n", port)
	server.Accept(lis)

}
