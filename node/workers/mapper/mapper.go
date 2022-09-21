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

func askForJoin(master string, client *rpc.Client) int {

	var hostname = os.Getenv("HOSTNAME")
	reqJoin := utils.JoinRequest{IP: hostname, Type: utils.MapperType}
	var reply_code int
	err := client.Call("Master.JoinMR", reqJoin, &reply_code)
	if err != nil {
		log.Fatal("Fatal error trying to join the cluster")
		os.Exit(-1)
	}
	return reply_code
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify master address:\n\t./mapper [addr:port]")
		os.Exit(1)
	}
	addr := os.Args[1]
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	reply := askForJoin(addr, client)

	if reply != 0 {
		log.Printf("Request declined from Master %s", addr)
		os.Exit(-1)
	} else {
		log.Printf("Request accepted from Master %s", addr)
	}

	worker := new(Worker)
	server := rpc.NewServer()

	err = server.Register(worker)
	if err != nil {
		log.Fatal("Error on Register(worker): ", err)
		os.Exit(-1)
	}

	var address string = ":" + strconv.Itoa(utils.WORKER_PORT)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error in listening:", err)
		os.Exit(-1)
	}
	log.Printf("Worker online on port %d\n", utils.WORKER_PORT)
	server.Accept(lis)

}
