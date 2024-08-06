/*
* @Author: chaidaxuan chaidaxuan@wps.cn
* @Date: 2024-07-26 16:58:52
* @LastEditors: chaidaxuan chaidaxuan@wps.cn
* @LastEditTime: 2024-07-26 17:40:37
* @FilePath: /urarao/GoProjects/Golang-Samples/src/Net/TCP_conn_init.go
* @Description:

	TCP连接的初始化示例

* Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
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
