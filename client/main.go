package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	workerBootingTime             = 100 * time.Millisecond
	connectionIdleTimeoutDuration = 10 * time.Second
)

var (
	workers      int
	queryStrings = []string{"test1", "test2", "test3", "test4", "quit"}
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.IntVar(&workers, "w", 100, "Concurrent client quantity")
	flag.Parse()
}

func main() {

	var wg sync.WaitGroup

	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go clientWorker(&wg, w)
		time.Sleep(workerBootingTime)
	}

	wg.Wait()
	log.Println("All worker finished")
}

func clientWorker(wg *sync.WaitGroup, w int) {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8000")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	conn.SetReadDeadline(time.Now().Add(connectionIdleTimeoutDuration))

	if err != nil {
		log.Printf("[Worker %d] %s\n", w, err.Error())
		return
	}
	defer conn.Close()
	log.Printf("[Work %d] connected!", w)

	for _, v := range queryStrings {
		log.Printf("[Worker %d] query %s\n", w, v)
		_, err := conn.Write([]byte(v + "\n"))
		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}

		responseMessage, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Reading string failed: %s\n", err.Error())
			continue
		}

		readLine := strings.TrimSuffix(responseMessage, "\n")

		switch readLine {
		case "quit":
			fmt.Println("quit")
		}

		time.Sleep(workerBootingTime)
	}

	wg.Done()
}
