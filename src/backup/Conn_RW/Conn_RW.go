package backup

import (
	"fmt"
	"net"
)

func client_test() {
	Tconn, _ := net.Dial("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", "8124"))
	buffer := make([]byte, 4096)
	sz, _ := Tconn.Read(buffer)
	if string(buffer[:sz-1]) == "123" {
		fmt.Println("123")
	}
	var s int
	fmt.Scanln(&s)
	Tconn.Close()
}

func server_test() {
	listener, _ := net.Listen("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", "8124"))
	defer listener.Close()
	// 与客户端建立Conn，初始化UserModel并转发给FileServer

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Server.Start, listener.Accept error: ", err)
			continue
		}
		sz, err := conn.Write([]byte("123" + "\n"))
		if err != nil {
			fmt.Println("ERROR,", err)
			return
		} else if sz == 0 {
			fmt.Println("sz == 0")
		}
	}
}
