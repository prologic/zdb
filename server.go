package main

import (
	// "bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"sync/atomic"

	zu "zdb/package/zulti"
)

var ops uint64
var gMapData *zu.ConcurrentMapHuge
var gDataname = "zdbdata"

func main() {
	gMapData = zu.NewConcurrentMapHuge()
	gMapData.Load(gDataname)

	// TIMER
	tickerMedi := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-tickerMedi.C:
				gMapData.Save(gDataname)
				fmt.Printf("PACK : %v\n", ops)
				fmt.Printf("DATA : %v\n", gMapData.Len())
				
			case <-quit:
				tickerMedi.Stop()
				return
			}
		}
	}()

	listener, err := net.Listen("tcp", "0.0.0.0:8989")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleClientRequest(con)
	}
}

func handleClientRequest(con net.Conn) {
	defer con.Close()

	_startTime := time.Now()
	startTime := _startTime.Unix()

	// scanner := bufio.NewScanner(con)
	// scanner.Scan()

	for i := 0; i < 10000000000; i++ {
		now := time.Now()
		sec := now.Unix()
		if i%1000000 == 0 {
			fmt.Printf("RPS: %v, %v, %v \n", sec, startTime, int64(i)/(sec-startTime+1))
		}

		gMapData.Set(uuid.New().String(), "What the hell")
		atomic.AddUint64(&ops, 1)
	}

	// for scanner.Scan() {
	// 	message := scanner.Text()
	// 	gMapData.Set(uuid.New().String(), message)
	// 	atomic.AddUint64(&ops, 1)
	// 	con.Write([]byte("GOT IT!\n"))
	// }
}
