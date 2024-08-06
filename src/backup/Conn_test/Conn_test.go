package backup

import (
	"fmt"
	"net"
)

func conninit() {
	// client
	serverIP, serverPort := "127.0.0.1", 8080
	conn1, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))

	// server
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	defer listener.Close()
	conn2, err := listener.Accept()
	if err != nil {
		fmt.Println("Server.Start, listener.Accept error: ", err)
	}

	// send
	chatmsg := "hello"
	_, err = conn1.Write([]byte(chatmsg))

	// receive
	buffer := make([]byte, 4096)
	sz, err := conn2.Read(buffer)
	if err != nil {
		fmt.Println("Server.HandleUser conn.Read error: ", err)
		return
	} else if sz == 0 {

	}
}
