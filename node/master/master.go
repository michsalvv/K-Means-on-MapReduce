package main

import (
	"kmeans-MR/utils"
	"log"
	"net"
	"net/rpc"
)

var cfg utils.Config

func main() {
	master := new(Master)
	cfg = utils.GetConfiguration()

	// Exposing RPCs
	server := rpc.NewServer()
	err := server.Register(master)
	if err != nil {
		log.Fatal("Error on register(master): ", err)
	}

	var address string = ":" + cfg.Server.MASTER_PORT
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error in listening:", err)
	}
	log.Printf("Master online on port [%s]\n\n", cfg.Server.MASTER_PORT)
	server.Accept(lis)
}
