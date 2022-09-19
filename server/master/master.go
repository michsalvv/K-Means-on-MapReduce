package main

import (
	"log"
	"net"
	"net/rpc"
	"strconv"
)

const MASTER_PORT int = 9001
const WORKER_IP string = "localhost"

func main() {

	master := new(Master)

	server := rpc.NewServer()
	err := server.Register(master)
	if err != nil {
		log.Fatal("Error on register(master): ", err)
	}

	var address string = ":" + strconv.Itoa(MASTER_PORT)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error in listening:", err)
	}
	log.Printf("Master online on port [%d]\n", MASTER_PORT)
	server.Accept(lis)
}
