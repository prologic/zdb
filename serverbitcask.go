package main

import (
	// "log"
	"time"
	"fmt"
	"strconv"

	"git.mills.io/prologic/bitcask"
	// "github.com/google/uuid"
)

func main() {
	_startTime := time.Now()
	startTime := _startTime.Unix()
	opts := []bitcask.Option{
		bitcask.WithMaxKeySize(1024),
		bitcask.WithMaxValueSize(4096),
	}

  db, _ := bitcask.Open("/tmp/db", opts...)
  defer db.Close()

  for i := 0; i < 100000000; i++ {
  	now := time.Now()
		sec := now.Unix()
		if i%10000 == 0 {
			fmt.Printf("RPS: %v, %v, %v \n", sec, startTime, int64(i)/(sec-startTime+1))
		}

  	db.Put([]byte(strconv.Itoa(i)), []byte("World"))
	}
  // val, _ := db.Get([]byte("Hello"))
  // log.Printf(string(val))
}
