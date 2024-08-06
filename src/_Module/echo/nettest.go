package nettest

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

var ip = "127.0.0.1"
var port = "8000"

func ServerRun() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer conn.Close()
		go func(conn net.Conn) {
			fmt.Println("new conn:", conn.RemoteAddr())
			lenbuf := make([]byte, 8)
			_, err := conn.Read(lenbuf)
			if err != nil {
				log.Println(err.Error())
				return
			}
			fmt.Println("lenbuf:", lenbuf)
			len := BytesToInt64(lenbuf)
			buffer := make([]byte, len)
			rsize, err := conn.Read(buffer)
			if err != nil {
				log.Println(err.Error())
				return
			}
			fmt.Println("buffer:", buffer)
			wsize, err := conn.Write(buffer)
			if err != nil {
				log.Println(err.Error())
				return
			}
			if wsize != rsize {
				log.Println("wsize != rsize", wsize, rsize)
			}
			log.Println("wsize == rsize", wsize, rsize)
			conn.Close()
		}(conn)
	}
}

func ClientRun() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer conn.Close()
	var str string
	_, err = fmt.Scan(&str)
	if err != nil {
		log.Println(err.Error())
		return
	}
	bstr := make([]byte, 0, len([]byte(str))+8)
	bstr = append(bstr, Int64ToBytes(int64(len([]byte(str))))...)
	bstr = append(bstr, []byte(str)...)
	fmt.Println("bstr:", bstr)
	wsize, err := conn.Write(bstr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	buffer := make([]byte, len([]byte(str)))
	time.Sleep(time.Second)
	rsize, err := conn.Read(buffer)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if wsize != rsize {
		log.Println("wsize != rsize", wsize, rsize)
	}
	fmt.Println(string(buffer))
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
