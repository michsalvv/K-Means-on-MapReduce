package main

import (
	"fmt"
	"kmeans-MR/utils"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cfg utils.Config

// This function will call JoinMR Master's RPC to ask for join the cluster
func askForJoin(master string, client *rpc.Client, workerType utils.WorkerType) int {

	var hostname = os.Getenv("HOSTNAME")
	reqJoin := utils.JoinRequest{IP: hostname, Type: workerType}
	var reply_code int
	err := client.Call("Master.JoinMR", reqJoin, &reply_code)
	if err != nil {
		log.Fatal("Fatal error trying to join the cluster", err.Error())
		os.Exit(-1)
	}
	return reply_code
}

func disconnect(master string, client *rpc.Client, workerType utils.WorkerType) {

	var reply int
	err := client.Call("Master.ExitMR", utils.JoinRequest{IP: os.Getenv("HOSTNAME"), Type: workerType}, &reply)
	if err != nil {
		log.Fatal("Fatal error trying to exit the cluster", err.Error())
		os.Exit(-1)
	}
}

func main() {
	cfg = utils.GetConfiguration()
	if len(os.Args) < 3 {
		fmt.Println("Please specify master address and type of tasker:\n\t./worker [addr:port] [reducer/mapper]")
		os.Exit(1)
	}
	// Connection to Master
	addr := os.Args[1]
	workerType := utils.DetectTaskType(os.Args[2])

	if workerType == -1 {
		log.Fatal("Please select a valid task type")
	}
	time.Sleep(3000) // wait master startup
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	// Handler og SIGINT signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		disconnect(addr, client, workerType)
		log.Printf("Disconnected")
		os.Exit(1)
	}()

	// Registration to Service
	for {
		log.Print("Connecting to master server...")
		reply := askForJoin(addr, client, workerType)
		if reply == 0 {
			log.Printf("Request accepted from Master %s", addr)
			break
		}
	}

	server := rpc.NewServer()
	if workerType == utils.Mapper {
		err = server.Register(new(Mapper))
		if err != nil {
			disconnect(addr, client, workerType)
			log.Fatal("Error on Register(worker): ", err)
			os.Exit(-1)
		}
	} else {
		err = server.Register(new(Reducer))
		if err != nil {
			disconnect(addr, client, workerType)
			log.Fatal("Error on Register(worker): ", err)
			os.Exit(-1)
		}
	}

	var address string = ":" + cfg.Server.WORKER_PORT
	lis, err := net.Listen("tcp", address)
	if err != nil {
		disconnect(addr, client, workerType)
		log.Fatal("Error in listening:", err)
		os.Exit(-1)
	}
	log.Printf("Worker online on port %s\n", cfg.Server.WORKER_PORT)
	server.Accept(lis)
}
