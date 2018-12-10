package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	connectionIdleTimeoutDuration = 10 * time.Second
	ratelimitPerSecond            = 30
	requestLimitRate              = time.Second / ratelimitPerSecond
)

var (
	mu                sync.RWMutex
	processedRequest  int
	connectingClients []string
	queryStrings      = make(chan string, 1024)
)

func main() {
	var tcpAddr *net.TCPAddr

	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8000")

	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	defer tcpListener.Close()

	go startStateServer()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		fmt.Println("A client connected : " + tcpConn.RemoteAddr().String())
		go tcpPipe(tcpConn)
	}

}

func tcpPipe(conn *net.TCPConn) {
	conn.SetReadDeadline(time.Now().Add(connectionIdleTimeoutDuration))
	ipString := conn.RemoteAddr().String()
	mu.Lock()
	connectingClients = append(connectingClients, ipString)
	mu.Unlock()

	defer func() {
		fmt.Println("disconnected :" + ipString)

		mu.Lock()
		for k, v := range connectingClients {
			if v == conn.RemoteAddr().String() {
				connectingClients = connectingClients[:k+copy(connectingClients[k:], connectingClients[k+1:])]
			}
		}
		mu.Unlock()
		conn.Close()
	}()
	reader := bufio.NewReader(conn)

	throttle := time.Tick(requestLimitRate)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		for i := 0; i < 30; i++ {
			go func() {
				for query := range queryStrings {
					<-throttle

					go func() {
						response := requestExternalAPI(query)
						if response.StatusCode == 200 {
							log.Printf("query %s success", query)
						} else {
							log.Printf("query %s fail", query)
						}
					}()

					mu.Lock()
					// log.Printf("Processed Request Count: %d", processedRequest)
					processedRequest++
					mu.Unlock()
				}
			}()
		}

		readLine := strings.TrimSuffix(message, "\n")

		fmt.Println(readLine)
		switch readLine {
		case "quit":
			fmt.Println("connection closing...")
			conn.Write([]byte("quit\n"))
			return
		default:
			conn.Write([]byte(readLine + "\n"))
			queryStrings <- readLine
		}

	}
}
