package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

const MASTER_IP string = "localhost"

type workerType int

const (
	MapperType  workerType = 1
	ReducerType workerType = 2
)

type JoinRequest struct {
	Port int
	Type workerType
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify port number:\n\t./reducer [port]")
		os.Exit(1)
	}
	addr := MASTER_IP + ":" + strconv.Itoa(9001)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	port, _ := strconv.Atoi(os.Args[1])

	var id int
	err = client.Call("Master.JoinGrep", JoinRequest{Port: port, Type: ReducerType}, &id)
	if err != nil {
		log.Fatal("Error in Master.JoinGrep: ", err.Error())
	}

	if id != 1 {
		log.Fatal("Request declined")
		os.Exit(1)
	} else {
		log.Println("Request accepted")
	}

	worker := new(Reducer)
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
