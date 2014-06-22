package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"proto"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	go func() {
		http.ListenAndServe(":9001", nil)
	}()
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			p := proto.NewProto(conn)
			for {
				cmd, data, err := p.ReadPackage()
				if err != nil {
					log.Println(err)
					return
				}
				err = p.WritePackage(cmd, data)
				if err != nil {
					log.Println(err)
					return
				}
				//log.Println("cmd: ", cmd, "data: ", string(data))
			}
		}()
	}
}
