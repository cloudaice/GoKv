package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"proto"
	"runtime"
	"sync"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	count int
	lock  *sync.Mutex
)

func Dial() {

	go func() {
		http.ListenAndServe(":9003", nil)
	}()
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
	}
	p := proto.NewProto(conn)
	for {
		err := p.WritePackage(0x55ff, []byte("hello world"))
		if err != nil {
			log.Fatal(err)
		}
		_, _, err = p.ReadPackage()
		if err != nil {
			log.Fatal(err)
		}
		lock.Lock()
		count++
		lock.Unlock()
	}

}

func main() {
	count = 0
	lock = new(sync.Mutex)
	for i := 0; i < 50; i++ {
		go Dial()
	}
	for {
		time.Sleep(time.Second)
		lock.Lock()
		log.Println(count)
		count = 0
		lock.Unlock()
	}
}
