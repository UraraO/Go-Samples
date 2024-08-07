/*=============
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-07 19:38:09
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-07 22:06:53
 FilePath: /Golang-Samples/src/_Utils/Net/http/httpServer.go
 Description:

 一个简单的http服务，使用golang标准库net/http

 Copyright (c) 2024 by UraraO, All Rights Reserved.
=============*/

package httptest

import (
	"fmt"
	"net"
	"net/http"
)

type ServerAddr struct {
	IP   string
	Port int
}

func (sa *ServerAddr) Address() string {
	return fmt.Sprintf("%s:%d", sa.IP, sa.Port)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello\n"))
}

func handleNull(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to UraraO's HTTP Server\n"))
}

func SimpleHttpServerTest(ip string, port int) error {
	addr := ServerAddr{
		IP:   ip,
		Port: port,
	}
	server := &http.Server{Addr: addr.Address(), Handler: nil}
	ln, err := net.Listen("tcp", addr.Address())
	if err != nil {
		return err
	}

	http.HandleFunc("GET /hello", handleHello)
	http.HandleFunc("GET /", handleNull)

	return server.Serve(ln)
}
