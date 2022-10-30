package main

import (
	"kmeans-MR/utils"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

func main() {
	master := new(Master)

	server := rpc.NewServer()
	err := server.Register(master)
	if err != nil {
		log.Fatal("Error on register(master): ", err)
	}

	var address string = ":" + strconv.Itoa(utils.MASTER_PORT)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error in listening:", err)
	}
	log.Printf("Master online on port [%d]\n\n", utils.MASTER_PORT)
	server.Accept(lis)
}
