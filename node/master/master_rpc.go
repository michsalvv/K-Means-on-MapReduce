package main

import (
	"bufio"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

const masterPath string = "server/master/"

var mappers []int
var reducers []int

type workerType int

const (
	MapperType  workerType = 1
	ReducerType workerType = 2
)

type JoinRequest struct {
	Port int
	Type workerType
}
type JoinResponse int
type Master int
type Input struct {
	Text       string
	WordToGrep string
}

func (m *Master) JoinGrep(req JoinRequest, reply *JoinResponse) error {

	switch req.Type {
	case MapperType:
		*reply, mappers = addWorker(req, mappers)
	case ReducerType:
		*reply, reducers = addWorker(req, reducers)
	}
	return nil
}

func addWorker(req JoinRequest, list []int) (JoinResponse, []int) {
	for _, x := range list {
		if x == int(req.Port) {
			return 0, list
		}
	}
	log.Printf("Request accepted from worker [%d]\n", req.Port)
	list = append(list, req.Port)
	return 1, list
}

func (m *Master) KMeans(in Input, reply *string) error {
	if len(mappers) == 0 {
		log.Fatal("Not possible permform grep: {0} mappers")
		return nil
	}
	log.Printf("Grepping word {%s} on file {%s}", in.WordToGrep, in.Text)

	file, err := os.Open(masterPath + in.Text)
	check(err)

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
	for i := 0; i < len(mappers); i++ {
		channels[i] = make(chan string)
		defer close(channels[i])
		go sendToMapper(splittedText[i], in.WordToGrep, mappers[i], channels[i])
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
func sendToMapper(lines, word string, port int, ch chan string) {
	addr := WORKER_IP + ":" + strconv.Itoa(port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing with worker: ", err)
	}
	defer client.Close()

	grepInput := Input{Text: lines, WordToGrep: word}
	var reply string
	err = client.Call("Worker.Map", grepInput, &reply)
	if err != nil {
		log.Fatal("Error in Worker.Map: ", err.Error())
	}

	ch <- reply

}

func sendToReduce(s string, port int, ch chan string) {
	addr := WORKER_IP + ":" + strconv.Itoa(port)
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

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}
