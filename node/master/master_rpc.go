package main

import (
	"bufio"
	"kmeans-MR/utils"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

// const masterPath string = "server/master/"

var mappers []utils.WorkerInfo
var reducers []utils.WorkerInfo

type Master int
type Input struct {
	Text       string
	WordToGrep string
}

func (m *Master) JoinMR(req utils.JoinRequest, reply *int) error {

	log.Printf("Request received from [%s] type [%d]", req.IP, req.Type)
	switch req.Type {

	case utils.Mapper:
		*reply, mappers = addWorker(req, mappers)

	case utils.Reducer:
		*reply, reducers = addWorker(req, reducers)
	}

	return nil
}

func addWorker(req utils.JoinRequest, workersList []utils.WorkerInfo) (int, []utils.WorkerInfo) {
	for _, x := range workersList {
		if x.IP == req.IP {
			log.Printf("Request from [%s] declined", req.IP)
			return -1, workersList
		}
	}
	log.Printf("Worker [%s] accepted in cluster\n", req.IP)
	workersList = append(workersList, utils.WorkerInfo{IP: req.IP})
	return 0, workersList
}

func removeWorker(req utils.JoinRequest, workersList []utils.WorkerInfo) (int, []utils.WorkerInfo) {
	var toRemove int
	for index, x := range workersList {
		if x.IP == req.IP {
			toRemove = index
			break
		}
	}
	workersList = append(workersList[:toRemove], workersList[toRemove+1:]...)
	return 0, workersList
}

func (m *Master) ExitMR(req utils.JoinRequest, reply *int) error {
	log.Printf("Worker [%s] disconnected", req.IP)
	switch req.Type {

	case utils.Mapper:
		*reply, mappers = removeWorker(req, mappers)

	case utils.Reducer:
		*reply, reducers = removeWorker(req, reducers)
	}
	return nil
}

// func (m *Master) KMeans(in Input, reply *string) error {
// 	log.Printf("MAPPERS ONLINE: %s", mappers)
// 	log.Printf("REDUCERS ONLINE: %s", reducers)
// 	return nil
// }

func (m *Master) KMeans(in Input, reply *string) error {
	if len(mappers) == 0 {
		log.Fatal("Not possible permform grep: {0} mappers")
		return nil
	}
	log.Printf("Grepping word {%s} on file {%s}", in.WordToGrep, in.Text)

	file, err := os.Open(in.Text)
	if err != nil {
		log.Print(err.Error())
		log.Print("ESCO?")
		return nil
	}

	reader := bufio.NewReader(file)
	var lines []string
	line, err := readLine(reader)
	for err == nil {
		if len(line) > 0 {
			lines = append(lines, line)
		}
		line, err = readLine(reader)
	}

	var splittedText = splitLines(lines, len(mappers))

	channels := make(map[int]chan string)
	for index, mapper := range mappers {
		channels[index] = make(chan string)
		defer close(channels[index])
		go sendToMapper(splittedText[index], in.WordToGrep, mapper, channels[index])
	}
	var response = sortResponse(channels)

	ch := make(chan string)
	defer close(ch)
	go sendToReduce(response, reducers[0], ch)
	*reply = <-ch
	return nil
}

func sortResponse(channels map[int]chan string) string {
	var s string
	// Waiting for #Workers replies
	for i := 0; i < len(channels); i++ {
		s += <-channels[i]
	}
	return s
}

func splitLines(lines []string, chunks int) []string {
	x := int(float64(len(lines)) / float64(chunks))
	log.Printf("[%d Lines]\t[%d Mappers]\t%d lines for mapper", len(lines), chunks, x)
	var offset, endOffset int
	splittedText := make([]string, chunks)
	for i := 0; i < chunks; i++ {
		offset = i * x
		if (i == chunks-1) && (chunks > 1) {
			endOffset = len(lines)
		} else {
			endOffset = offset + x
		}
		log.Printf("mapper: %d\tfrom: %d \tto: %d\n", i, offset, endOffset)
		for j := offset; j < endOffset; j++ {
			splittedText[i] += lines[j] + "\n"
		}
		splittedText[i] = strings.TrimSuffix(splittedText[i], "\n") // Delete last \n character
	}
	return splittedText
}
func sendToMapper(lines, word string, mapper utils.WorkerInfo, ch chan string) {
	addr := mapper.IP + ":" + strconv.Itoa(utils.WORKER_PORT)
	log.Print("dialing with ", addr)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing with worker: ", err)
	}
	defer client.Close()

	grepInput := Input{Text: lines, WordToGrep: word}
	var reply string
	err = client.Call("Mapper.Map", grepInput, &reply)
	if err != nil {
		log.Fatal("Error in Mapper.Map: ", err.Error())
	}

	ch <- reply

}

func sendToReduce(s string, reducer utils.WorkerInfo, ch chan string) {
	addr := reducer.IP + ":" + strconv.Itoa(utils.WORKER_PORT)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing with worker: ", err)
	}
	defer client.Close()

	var reply string
	err = client.Call("Reducer.Reduce", s, &reply)
	if err != nil {
		log.Fatal("Error in Reducer.Reduce: ", err.Error())
	}

	ch <- reply

}

// Readln returns a single line (without the ending \n) from the input buffered reader.
func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
