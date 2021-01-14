package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"sync/atomic"
)

var ops uint64

func main() {
	// TIMER
	tickerMedi := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-tickerMedi.C:
				fmt.Printf("PACK : %v\n", ops)
				// fmt.Printf("DATA : %v\n", gMapData.Len())
			case <-quit:
				tickerMedi.Stop()
				return
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	con, err := net.Dial("tcp", "0.0.0.0:8989")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	// Receiver
	go clientReceiver(con)

	// Sender
	go clientSender(con)

	// Waiting receiver done
	wg.Wait()
}

func clientReceiver(con net.Conn) {
	scanner := bufio.NewScanner(con)
	for scanner.Scan() {
		_ = scanner.Text()
		atomic.AddUint64(&ops, 1)
	}
}

func clientSender(con net.Conn) {
	_startTime := time.Now()
	startTime := _startTime.Unix()

	for i := 0; i < 10000000; i++ {
		now := time.Now()
		sec := now.Unix()
		if i%10000 == 0 {
			fmt.Printf("RPS: %v, %v, %v \n", sec, startTime, int64(i)/(sec-startTime+1))
		}

		message := strings.TrimSpace("Hello server !!!")
		con.Write([]byte(message + "\n"))
	}
}
